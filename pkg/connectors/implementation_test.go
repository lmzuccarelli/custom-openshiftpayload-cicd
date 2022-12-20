package connectors

import (
	"os"
	"testing"

	"github.com/microlib/logger/pkg/multi"
)

var workingDir string

func TestAllConnectors(t *testing.T) {
	logger := multi.NewLogger(multi.COLOR, "trace")
	client := NewClientConnections(logger)
	workingDir, _ = os.Getwd()
	client.Info("test %s", "simple")
	client.Debug("test %s", "simple")
	client.Trace("test %s", "simple")
	client.Error("test %s", "simple")
	t.Run("Testing ExecutePipeline : should pass", func(t *testing.T) {
		err := client.ExecOS(".", "ls", []string{"-la"}, true)
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})
}
