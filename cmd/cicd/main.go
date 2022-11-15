package main

import (
	"flag"
	"os"

	"github.com/microlib/logger/pkg/multi"

	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/connectors"
	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/service"
)

var (
	logLevel             string
	pipelineFile         string
	configFile           string
	generateTaskRunFiles string
	help                 string
)

func init() {
	flag.StringVar(&logLevel, "l", "info", "Set log level [info,debug,trace]")
	flag.StringVar(&pipelineFile, "p", "", "Path and name of the pipeline yaml file")
	flag.StringVar(&configFile, "c", "", "Use config file - this overrides the taskruns to execute")
	flag.StringVar(&generateTaskRunFiles, "g", "", "Use this flag to generate all runtasks from a given text file")
	flag.StringVar(&help, "h", " ", "Display usage")
}

func main() {

	flag.Parse()
	if help == "" {
		flag.Usage()
		os.Exit(1)
	}

	logger := multi.NewLogger(multi.COLOR, logLevel)
	client := connectors.NewClientConnections(logger)

	if len(generateTaskRunFiles) > 0 {
		err := service.GenerateTaskRunFiles(generateTaskRunFiles)
		if err != nil {
			client.Error("generating taskrun files %v", err)
			os.Exit(1)
		}
		client.Info("completed generating taskrun files")
		os.Exit(0)
	}

	err := service.ExecutePipeline(pipelineFile, configFile, client)
	if err != nil {
		client.Error("pipeline execution %v", err)
		os.Exit(1)
	}
	client.Info("completed executing pipeline")
}
