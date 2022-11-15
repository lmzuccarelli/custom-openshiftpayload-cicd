package service

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/schema"
	"gopkg.in/yaml.v2"
)

var taskRunTemplate = `
Kind: TaskRun
apiVersion: 0.0.1
metadata:
  name: {{ .Name }}
spec:
  taskRef:  {{ .TaskRef }}
  parameters:
    - name: repository-url
      value: '{{ .GitURL }}'
    - name: repository-name
      value: '{{ .Name }}'
    - name: branch
      value: 'master'
    - name: dockerfile
      value: '{{ .Dockerfile }}'
`

func deepCopyTask(task schema.Task) schema.Task {
	t := schema.Task{}
	t.Kind = task.Kind
	t.APIVersion = task.APIVersion
	t.Metadata.Name = task.Metadata.Name
	s := schema.Step{}
	for _, step := range task.Spec.Steps {
		s.Workspace = step.Workspace
		cmd := []string{}
		args := []string{}
		s.Name = step.Name
		s.Description = step.Description
		for _, s := range step.Command {
			cmd = append(cmd, s)
		}
		for _, a := range step.Args {
			args = append(args, a)
		}
		s.Command = cmd
		s.Args = args
		t.Spec.Steps = append(t.Spec.Steps, s)
	}
	return t
}

func mergeParams(tr *schema.TaskRun, pt *schema.PipelineTemplate) {
	for _, p := range pt.Spec.Parameters {
		tr.Spec.Parameters = append(tr.Spec.Parameters, p)
	}
}

func updateParameters(task schema.Task, taskrun schema.TaskRun) {
	for x, step := range task.Spec.Steps {
		for y, arg := range step.Args {
			if strings.Contains(arg, "${params.") {
				task.Spec.Steps[x].Args[y] = replaceWithParamValue(taskrun, arg)
			}
		}
		if strings.Contains(step.Workspace, "${params.") {
			task.Spec.Steps[x].Workspace = replaceWithParamValue(taskrun, step.Workspace)
		}
	}
}

func replaceWithParamValue(task schema.TaskRun, value string) string {
	newValue := value
	for _, param := range task.Spec.Parameters {
		newValue = strings.ReplaceAll(newValue, "${params."+param.Name+"}", param.Value)
	}
	return newValue
}

func findRelatedTask(tasks []schema.Task, reference string) schema.Task {
	for _, task := range tasks {
		if reference == task.Metadata.Name {
			return task
		}
	}
	return schema.Task{}
}

func findTaskRun(taskruns []schema.TaskRun, reference string) schema.TaskRun {
	for _, taskrun := range taskruns {
		if taskrun.Metadata.Name == reference {
			return taskrun
		}
	}
	return schema.TaskRun{}
}

func readAllTaskFiles(dir string) ([]schema.Task, error) {
	var tasks []schema.Task
	var task *schema.Task
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []schema.Task{}, err
	}
	for _, f := range files {
		yfile, err := ioutil.ReadFile(dir + f.Name())
		if err != nil {
			return []schema.Task{}, err
		}
		err = yaml.Unmarshal(yfile, &task)
		if err != nil {
			return []schema.Task{}, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func readAllTaskRunFiles(config, dir string) ([]schema.TaskRun, error) {
	var taskruns []schema.TaskRun
	var taskrun *schema.TaskRun
	var taskrunConfig *schema.TaskRunConfig

	// read the config file - ignore errors
	if len(config) > 0 {
		cfg, err := ioutil.ReadFile(config)
		if err != nil {
			return []schema.TaskRun{}, err
		}
		err = yaml.Unmarshal(cfg, &taskrunConfig)
		if err != nil {
			return []schema.TaskRun{}, err
		}
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []schema.TaskRun{}, err
	}
	for _, f := range files {
		yfile, err := ioutil.ReadFile(dir + f.Name())
		if err != nil {
			return []schema.TaskRun{}, err
		}
		err = yaml.Unmarshal(yfile, &taskrun)
		if err != nil {
			return []schema.TaskRun{}, err
		}
		taskruns = append(taskruns, *taskrun)
	}

	// filter if config is set
	if len(config) > 0 {
		var tr []schema.TaskRun
		for _, cfg := range taskrunConfig.Spec.TaskRunList {
			tr = append(tr, findTaskRun(taskruns, cfg))
		}
		return tr, nil
	}
	return taskruns, nil
}

func readPipelineTemplateFile(file string) (*schema.PipelineTemplate, error) {
	var pt *schema.PipelineTemplate
	yfile, err := ioutil.ReadFile(file)
	if err != nil {
		return &schema.PipelineTemplate{}, err
	}
	err = yaml.Unmarshal(yfile, &pt)
	if err != nil {
		return &schema.PipelineTemplate{}, err
	}
	return pt, nil
}

func generateTaskRunFiles(file string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(f), "\n")
	for _, line := range lines {
		hld := strings.Split(line, "/")
		schema := &schema.TemplateSchema{Name: hld[len(hld)-1], GitURL: line, Dockerfile: "Dockerfile.ocp", TaskRef: "custom-openshiftbuild-task"}
		//parse some content and generate a template
		tmpl := template.New("taskruns")
		//parse some content and generate a template
		tmp, _ := tmpl.Parse(taskRunTemplate)
		var tpl bytes.Buffer
		tmp.Execute(&tpl, schema)
		err = ioutil.WriteFile("./manifests/taskruns/"+schema.Name+".yaml", tpl.Bytes(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
