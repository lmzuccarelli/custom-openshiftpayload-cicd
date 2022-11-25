package service

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/lmzuccarelli/custom-tekton-emulator-cicd/pkg/connectors"
	"github.com/lmzuccarelli/custom-tekton-emulator-cicd/pkg/schema"
)

func ExecutePipeline(path string, c connectors.Clients) error {

	var p schema.Pipeline
	var t []schema.Task
	var tr []schema.TaskRun
	var ok bool

	os.Chdir(path)
	// read the main kustomization yaml file
	hldK, err := yamlToStruct("./kustomization.yaml", schema.Kustomization{})
	k, ok := hldK.(schema.Kustomization)
	if err != nil && !ok {
		return err
	}
	c.Trace("main kustomization file : %v", k)

	// get all the relevant files to load
	for _, f := range k.Bases {
		os.Chdir(f)
		subK, err := yamlToStruct("./kustomization.yaml", schema.Kustomization{})
		k, ok := subK.(schema.Kustomization)
		if err != nil && !ok {
			return err
		}
		c.Trace("sub level kustomization file %s : %v", f, subK)

		switch {
		case strings.Contains(f, "pipelines"):
			// only the first pipeline file will be used
			// other pipleine files will be ignored
			hldPt, err := yamlToStruct(f+"/"+k.Bases[0], schema.Pipeline{})
			p, ok = hldPt.(schema.Pipeline)
			if err != nil && !ok {
				return err
			}
			c.Debug("pipelines : %v", p)
			if len(p.Spec.Workspaces) == 0 {
				c.Info("field 'Workspaces is empty current directory 'working-dir' is set as default")
				p.Spec.Workspaces = append(p.Spec.Workspaces, schema.Workspace{Name: "./working-dir"})
			}
			c.Info("deleting working directory %s", p.Spec.Workspaces[0].Name)
			err = os.RemoveAll(p.Spec.Workspaces[0].Name)
			if err != nil {
				return err
			}
			c.Info("creating working directory %s", p.Spec.Workspaces[0].Name)
			err = os.MkdirAll(p.Spec.Workspaces[0].Name, os.ModePerm)
			if err != nil {
				return err
			}
		case strings.Contains(f, "tasks"):
			t, err = readAllTaskFiles(f, k.Bases)
			if err != nil {
				return err
			}
			c.Debug("tasks : %v", t)
		case strings.Contains(f, "taskruns"):
			tr, err = readAllTaskRunFiles(f, k.Bases)
			if err != nil {
				return err
			}
			c.Debug("taskruns : %v", tr)
			if len(tr) == 0 {
				return fmt.Errorf("taskrun list is empty (verify taskruns kustomization file)")
			}
		}
	}

	// we now execute the pipeline from the taskruns
	for _, taskrun := range tr {
		mergeParams(&taskrun, &p)
		task := findRelatedTask(t, taskrun.Spec.TaskRef.Name)
		if reflect.DeepEqual(task, schema.Task{}) {
			return fmt.Errorf("no related task '%s' found to execute", taskrun.Spec.TaskRef.Name)
		}
		newTask := deepCopyTask(task)
		updateParameters(newTask, taskrun)
		for _, step := range newTask.Spec.Steps {
			c.Info("executing %s for %s", step.Name, taskrun.Metadata.Name)
			err := c.ExecOS(p.Spec.Workspaces[0].Name+"/"+step.Workspace, step.Command[0], step.Args, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GenerateTaskRunFiles(path, dstPath string) error {
	bcs, err := readAllBuildConfigs(path)
	if err != nil {
		return err
	}
	return generateTaskRunFiles(dstPath, bcs)
}
