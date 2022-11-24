package service

import (
	"testing"

	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/connectors"
	"github.com/microlib/logger/pkg/multi"
)

func TestExecutePipeline(t *testing.T) {
	logger := multi.NewLogger(multi.COLOR, "trace")
	client := connectors.NewClientConnections(logger)

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
	t.Run("Testing GenerateTaskRunFiles : should pass", func(t *testing.T) {
		err := GenerateTaskRunFiles("../../tests/buildconfigs", "../../tests/taskruns")
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})
}
