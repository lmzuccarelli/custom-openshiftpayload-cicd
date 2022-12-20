package main

import (
	"flag"
	"os"

	"github.com/microlib/logger/pkg/multi"

	"github.com/lmzuccarelli/custom-tekton-emulator-cicd/pkg/connectors"
	"github.com/lmzuccarelli/custom-tekton-emulator-cicd/pkg/service"
)

var (
	logLevel             string
	kustomizePath        string
	configFile           string
	buildConfigDir       string
	generateTaskRunFiles string
	help                 string
)

func init() {
	flag.StringVar(&logLevel, "l", "info", "Set log level [info,debug,trace]")
	flag.StringVar(&kustomizePath, "k", "", "Path for the initial kustomization file")
	flag.StringVar(&configFile, "c", "", "Use config file - this overrides the taskruns to execute")
	flag.StringVar(&generateTaskRunFiles, "g", "", "Use this flag to set the destination folder for taskrun object to be saved to")
	flag.StringVar(&buildConfigDir, "b", "", "Use this flag to generate all runtasks from a given buildconfigs directory")
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

	if len(generateTaskRunFiles) > 0 && len(buildConfigDir) > 0 {
		err := service.GenerateTaskRunFiles(buildConfigDir, generateTaskRunFiles)
		if err != nil {
			client.Error("generating taskrun files %v", err)
			os.Exit(1)
		}
		client.Info("completed generating taskrun files")
		os.Exit(0)
	}

	if len(kustomizePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	err := service.ExecutePipeline(kustomizePath, client)
	if err != nil {
		client.Error("pipeline execution %v", err)
		os.Exit(1)
	}
	client.Info("completed executing pipeline")
}
