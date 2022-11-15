package connectors

import (
	"os"
	"os/exec"

	"github.com/microlib/logger/pkg/multi"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	Logger *multi.Logger
}

func NewClientConnections(logger *multi.Logger) Clients {
	return &Connectors{Logger: logger}
}

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Errorf(msg, val...)
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Infof(msg, val...)
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debugf(msg, val...)
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Tracef(msg, val...)
}

func (c *Connectors) ExecOS(path string, command string, params []string, trim bool) error {
	cmd := exec.Command(command, params...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
