package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/connectors"
)

func ExecutePipeline(file, config string, c connectors.Clients) error {

	pt, err := readPipelineTemplateFile(file)
	if err != nil {
		return err
	}
	if len(pt.Spec.WorkingDirectory) == 0 {
		c.Info("field 'workingDirectory is empty current directory 'working-dir' is set as default")
		pt.Spec.WorkingDirectory = "./working-dir"
	}
	c.Info("deleting working directory %s", pt.Spec.WorkingDirectory)
	os.RemoveAll(pt.Spec.WorkingDirectory)
	c.Info("creating working directory %s", pt.Spec.WorkingDirectory)
	os.MkdirAll(pt.Spec.WorkingDirectory, os.ModePerm)

	dir := "./" + strings.Split(filepath.Dir(file), "/")[0]

	tasks, err := readAllTaskFiles(dir + "/tasks/")
	if err != nil {
		return err
	}

	taskruns, err := readAllTaskRunFiles(config, dir+"/taskruns/")
	if err != nil {
		return err
	}

	for _, taskrun := range taskruns {
		mergeParams(&taskrun, pt)
		task := findRelatedTask(tasks, taskrun.Spec.TaskRef)
		newTask := deepCopyTask(task)
		updateParameters(newTask, taskrun)
		for _, step := range newTask.Spec.Steps {
			c.Info("executing %s for %s", step.Name, taskrun.Metadata.Name)
			err := c.ExecOS(pt.Spec.WorkingDirectory+"/"+step.Workspace, step.Command[0], step.Args, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GenerateTaskRunFiles(file string) error {
	return generateTaskRunFiles(file)
}
