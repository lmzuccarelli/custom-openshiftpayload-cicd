package service

import (
	"testing"

	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/connectors"
	"github.com/microlib/logger/pkg/multi"
)

func TestExecutePipeline(t *testing.T) {

	logger := multi.NewLogger(multi.COLOR, "debug")
	client := connectors.NewClientConnections(logger)

	t.Run("Testing Convert : should pass", func(t *testing.T) {
		err := ExecutePipeline("../../tests/pipeline/pipeline.yaml", "../../tests/config/config.yaml", client)
		if err != nil {
			t.Fatalf("Should not fail : found error %v", err)
		}
	})
	/*
		t.Run("Testing Convert : (bad config) should fail", func(t *testing.T) {
			err := Convert("../../tests/conf.yaml")
			if err == nil { e
				t.Fatalf("Should fail : found error %v", err)
			}
		})

		t.Run("Testing Convert : (bad bc) should fail", func(t *testing.T) {
			err := Convert("../../tests/config-bad-bc.yaml")
			if err == nil {
				t.Fatalf("Should fail : found error %v", err)
			}
		})
	*/
}
