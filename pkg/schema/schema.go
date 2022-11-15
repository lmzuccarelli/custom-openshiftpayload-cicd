package schema

type Task struct {
	Kind       string `yaml:"Kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Steps []Step
	} `yaml:"spec"`
}

type Step struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Command     []string `yaml:"command"`
	Args        []string `yaml:"args"`
	Workspace   string   `yaml:"workspace,omitempty"`
}

type TaskRun struct {
	Kind       string `yaml:"Kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		TaskRef    string `yaml:"taskRef"`
		Parameters []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"parameters"`
	} `yaml:"spec"`
}

type PipelineTemplate struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		WorkingDirectory string `yaml:"workingDirectory"`
		Parameters       []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"parameters"`
		Tasks []struct {
			Name    string `yaml:"name"`
			TaskRef string `yaml:"taskRef"`
		} `yaml:"tasks"`
	} `yaml:"spec"`
}

type TemplateSchema struct {
	Name       string
	GitURL     string
	Dockerfile string
	TaskRef    string
}

type TaskRunConfig struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		TaskRunList []string `yaml:"taskRunList"`
	} `yaml:"spec"`
}
