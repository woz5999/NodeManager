package kubectl

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Kubectl struct{}

func (k Kubectl) Exec(args []string) error {
	cmd := exec.Command("kubectl", args...)
	log.Info("Running command ")
	err := cmd.Run()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
