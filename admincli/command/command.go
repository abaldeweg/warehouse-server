package command

import (
	"os/exec"
)

func Command(cmd []string) ([]byte, error) {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
