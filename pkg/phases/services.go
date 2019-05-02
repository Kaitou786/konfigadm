package phases

import "fmt"

import . "github.com/moshloop/configadm/pkg/types"

var Services Phase = services{}

type services struct{}

func (p services) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	for name, svc := range sys.Services {
		filename := fmt.Sprintf("/etc/systemd/system/%s.service", name)
		files[filename] = File{Content: svc.Extra.ToUnitFile()}
		commands = append(commands, Command{Cmd: "systemctl enable " + name})
		commands = append(commands, Command{Cmd: "systemctl start " + name})
	}
	return commands, files, nil
}