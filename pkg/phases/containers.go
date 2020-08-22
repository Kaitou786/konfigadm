package phases

import (
	"fmt"
	"os"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
)

var Containers Phase = containers{}

type containers struct{}

func (p containers) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}
	for _, c := range sys.Containers {

		sys.Services[c.Name()] = Service{
			Name:      c.Name(),
			ExecStart: exec(sys, c),
			Extra:     DefaultSystemdService(c.Name()),
		}
		if len(c.Env) > 0 {
			files["/etc/environment."+c.Name()] = File{Content: toEnvironmentFile(ctx, c)}
		}

	}
	return commands, files, nil
}

func (p containers) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	for f := range cfg.Files {

		if _, err := os.Stat(f); err != nil {
			verify = false
			results.Fail("%s does not exist", f)
		} else {
			results.Pass("%s exists", f)
		}
	}

	for f := range cfg.Templates {
		if _, err := os.Stat(f); err != nil {
			verify = false
			results.Fail("%s does not exist", f)
		} else {
			results.Pass("%s exists", f)
		}
	}

	return verify
}

func toEnvironmentFile(ctx *SystemContext, c Container) string {
	s := ""
	for k, v := range c.Env {
		s += fmt.Sprintf("%s=%s\n", k, v)
	}
	return s
}

func exec(sys *Config, c Container) string {
	exec := c.DockerOpts
	if len(c.Env) > 0 {
		exec += fmt.Sprintf(" --env-file /etc/environment.%s", c.Name())
	}
	if c.Network != "" {
		exec += " --network " + c.Network
	}

	for _, v := range c.Volumes {
		exec += fmt.Sprintf(" -v %s", v)
	}

	for _, p := range c.Ports {
		exec += fmt.Sprintf(" -p %d:%d", p.Port, p.Target)
	}

	return fmt.Sprintf("%s run --rm --name %s %s %s %s", sys.ContainerRuntime.GetCLI(), c.Name(), exec, c.Image, c.Args)
}
