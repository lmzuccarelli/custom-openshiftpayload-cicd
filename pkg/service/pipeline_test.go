package service

import (
	"os"
	"testing"

	"github.com/lmzuccarelli/custom-tekton-emulator-cicd/pkg/connectors"
	"github.com/microlib/logger/pkg/multi"
)

var workingDir string

func TestExecutePipeline(t *testing.T) {
	logger := multi.NewLogger(multi.COLOR, "trace")
	client := connectors.NewClientConnections(logger)
	workingDir, _ = os.Getwd()
	t.Run("Testing ExecutePipeline : should pass", func(t *testing.T) {
		err := ExecutePipeline("../../tests/config", client)
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})

	t.Run("Testing ExecutePipeline : should pass", func(t *testing.T) {
		err := ExecutePipeline("../../tests/config-bad", client)
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})

}

func TestGenerateTaskRunFiles(t *testing.T) {
	os.Chdir(workingDir)
	t.Run("Testing GenerateTaskRunFiles : should pass", func(t *testing.T) {
		err := GenerateTaskRunFiles("../../tests/buildconfigs", "../../tests/taskruns")
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})
}
