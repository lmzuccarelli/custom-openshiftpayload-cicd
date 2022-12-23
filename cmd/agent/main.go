package main

import (
	"flag"
	"os"

	"github.com/lmzuccarelli/custom-tekton-emulator-cicd/pkg/connectors"
	"github.com/microlib/logger/pkg/multi"
)

var (
	logLevel       string
	executeTaskRun string
	help           string
)

func init() {
	flag.StringVar(&logLevel, "l", "info", "Set log level [info,debug,trace]")
	flag.StringVar(&executeTaskRun, "t", "", "Use this flag to execute a specific task run")
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
	basePath, err := os.Getwd()
	if err != nil {
		client.Error("executing agent %v", err)
		os.Exit(1)
	}

	err = os.Remove(basePath + "/logs/" + executeTaskRun + ".log")
	if err != nil {
		client.Error("executing agent %v", err)
	}

	err = os.WriteFile(basePath+"/logs/"+executeTaskRun+".log", []byte(""), 0755)
	if err != nil {
		client.Error("executing agent %v", err)
		os.Exit(1)
	}
	err = client.ExecOS(".", "./build/cicd", []string{"-k", "environments/overlays/cicd", "-l", logLevel, "-t", executeTaskRun}, basePath+"/logs/"+executeTaskRun+".log")
	if err != nil {
		client.Error("executing agent %v", err)
		os.Exit(1)
	}
	client.Info("completed agent execution")
}
