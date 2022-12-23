package connectors

import (
	"bytes"
	"io"
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

func (c *Connectors) ExecOS(path string, command string, params []string, logFile string) error {
	var out bytes.Buffer
	multi := io.MultiWriter(os.Stdout, &out)
	cmd := exec.Command(command, params...)
	cmd.Dir = path
	cmd.Stdout = multi
	cmd.Stderr = multi

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if len(logFile) > 0 {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			return err
		}

		defer file.Close()
		_, err = file.WriteString(out.String())
		if err != nil {
			return err
		}
	}

	return nil
}
