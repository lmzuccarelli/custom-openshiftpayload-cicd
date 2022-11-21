package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/luigizuccarelli/custom-openshiftpayload-cicd/pkg/schema"
	"gopkg.in/yaml.v2"
)

const (
	build           string = "Build"
	buildConfig     string = "BuildConfig"
	manifests       string = "/manifests/"
	dotYml          string = ".yaml"
	errMsgYaml      string = "converting struct to yaml for %s"
	errMsgUnmarshal string = "unmarshaling yaml to struct %v for %s"
)

var taskRunTemplate = `
kind: TaskRun
apiVersion: 0.0.1
metadata:
  name: {{ .Name }}
spec:
  taskRef:
    name: {{ .TaskRef }}
  params:
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

func mergeParams(tr *schema.TaskRun, pt *schema.Pipeline) {
	for _, p := range pt.Spec.Params {
		tr.Spec.Params = append(tr.Spec.Params, schema.Param{Name: p.Name, Value: p.Default})
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
	for _, param := range task.Spec.Params {
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

func readAllTaskFiles(dir string, files []string) ([]schema.Task, error) {
	var tasks []schema.Task
	for _, f := range files {
		hldTask, err := yamlToStruct(dir+"/"+f, schema.Task{})
		t, ok := hldTask.(schema.Task)
		if err != nil && ok {
			return []schema.Task{}, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func readAllTaskRunFiles(dir string, files []string) ([]schema.TaskRun, error) {
	var taskruns []schema.TaskRun
	var tr schema.TaskRun
	var ok bool

	for _, f := range files {
		hldTr, err := yamlToStruct(dir+"/"+f, schema.TaskRun{})
		tr, ok = hldTr.(schema.TaskRun)
		if err != nil && !ok {
			return []schema.TaskRun{}, err
		}
		taskruns = append(taskruns, tr)
	}
	return taskruns, nil
}

// readAllBuildConfigs - reads all the BuildConfigs from a given directory
func readAllBuildConfigs(dir string) ([]schema.BuildConfig, error) {
	var bcs []schema.BuildConfig
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []schema.BuildConfig{}, err
	}
	for _, f := range files {
		if !f.IsDir() {
			hldBC, err := yamlToStruct(dir+"/"+f.Name(), schema.BuildConfig{})
			bc, ok := hldBC.(schema.BuildConfig)
			if err != nil && ok {
				return []schema.BuildConfig{}, err
			}
			// we are only interested in BuildCobcnfigs
			if bc.Kind == buildConfig {
				bcs = append(bcs, bc)
			}
		}
	}
	return bcs, nil
}

func generateTaskRunFiles(bcs []schema.BuildConfig) error {
	for _, bc := range bcs {
		schema := &schema.TemplateSchema{Name: filepath.Base(bc.Spec.Source.Git.URI), GitURL: bc.Spec.Source.Git.URI, Dockerfile: bc.Spec.Strategy.DockerStrategy.DockerfilePath, TaskRef: "custom-openshift-build"}
		//parse some content and generate a template
		tmpl := template.New("taskruns")
		tmp, _ := tmpl.Parse(taskRunTemplate)
		var tpl bytes.Buffer
		tmp.Execute(&tpl, schema)
		err := ioutil.WriteFile("./manifests/taskruns/"+schema.Name+".yaml", tpl.Bytes(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func yamlToStruct(file string, strctType interface{}) (interface{}, error) {
	yfile, err := ioutil.ReadFile(file)
	if err != nil {
		return strctType, err
	}
	v := reflect.TypeOf(strctType).String()
	switch v {
	case "schema.BuildConfig":
		var bc schema.BuildConfig
		err = yaml.Unmarshal(yfile, &bc)
		if err != nil {
			return strctType, fmt.Errorf(errMsgUnmarshal, err, file)
		}
		return bc, nil
	case "schema.TaskRun":
		var tr schema.TaskRun
		err = yaml.Unmarshal(yfile, &tr)
		if err != nil {
			return strctType, fmt.Errorf(errMsgUnmarshal, err, file)
		}
		return tr, nil
	case "schema.Task":
		var t schema.Task
		err = yaml.Unmarshal(yfile, &t)
		if err != nil {
			return strctType, fmt.Errorf(errMsgUnmarshal, err, file)
		}
		return t, nil
	case "schema.Pipeline":
		var p schema.Pipeline
		err = yaml.Unmarshal(yfile, &p)
		if err != nil {
			return strctType, fmt.Errorf(errMsgUnmarshal, err, file)
		}
		return p, nil
	case "schema.Kustomization":
		var k schema.Kustomization
		err = yaml.Unmarshal(yfile, &k)
		if err != nil {
			return strctType, fmt.Errorf(errMsgUnmarshal, err, file)
		}
		return k, nil
	case "schema.TaskRunConfig":
		var trc schema.TaskRunConfig
		err = yaml.Unmarshal(yfile, &trc)
		if err != nil {
			return strctType, fmt.Errorf(errMsgUnmarshal, err, file)
		}
		return trc, nil
	}
	return strctType, nil
}
