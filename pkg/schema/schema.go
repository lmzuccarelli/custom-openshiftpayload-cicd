package schema

type Param struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

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
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Params    []Param `yaml:"params"`
		Resources struct {
		} `yaml:"resources"`
		ServiceAccountName string `yaml:"serviceAccountName"`
		TaskRef            struct {
			Kind string `yaml:"kind"`
			Name string `yaml:"name"`
		} `yaml:"taskRef"`
		Timeout    string `yaml:"timeout"`
		Workspaces []struct {
			Name                  string `yaml:"name"`
			PersistentVolumeClaim struct {
				ClaimName string `yaml:"claimName"`
			} `yaml:"persistentVolumeClaim"`
		} `yaml:"workspaces"`
	} `yaml:"spec"`
}

type Workspace struct {
	Name string `yaml:"name"`
}

type Pipeline struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Params []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Default     string `yaml:"default,omitempty"`
			Type        string `yaml:"type"`
		} `yaml:"params"`
		Workspaces []Workspace `yaml:"workspaces"`
		Resources  []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
		} `yaml:"resources"`
		Tasks []struct {
			Name    string  `yaml:"name"`
			Params  []Param `yaml:"params"`
			TaskRef struct {
				Kind string `yaml:"kind"`
				Name string `yaml:"name"`
			} `yaml:"taskRef,omitempty"`
			Workspaces []struct {
				Name      string `yaml:"name"`
				Workspace string `yaml:"workspace"`
			} `yaml:"workspaces"`
			RunAfter []string `yaml:"runAfter,omitempty"`
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

type BuildConfig struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		RunPolicy string `yaml:"runPolicy"`
		Triggers  []struct {
			Type   string `yaml:"type"`
			Github struct {
				Secret string `yaml:"secret"`
			} `yaml:"github,omitempty"`
			Generic struct {
				Secret string `yaml:"secret"`
			} `yaml:"generic,omitempty"`
		} `yaml:"triggers"`
		Source struct {
			Git struct {
				URI string `yaml:"uri"`
			} `yaml:"git"`
			ContextDir string `yaml:"contextDir"`
			Images     []struct {
				From struct {
					Kind string `yaml:"kind"`
					Name string `yaml:"name"`
				} `yaml:"from"`
				Paths []struct {
					SourcePath     string `yaml:"sourcePath"`
					DestinationDir string `yaml:"destinationDir"`
				} `yaml:"paths"`
			} `yaml:"images"`
		} `yaml:"source"`
		Strategy struct {
			Type           string `yaml:"type"`
			SourceStrategy struct {
				From struct {
					Kind string `yaml:"kind"`
					Name string `yaml:"name"`
				} `yaml:"from"`
			} `yaml:"sourceStrategy"`
			DockerStrategy struct {
				ImageOptimizationPolicy string `yaml:"imageOptimizationPolicy"`
				DockerfilePath          string `yaml:"dockerfilePath"`
				From                    struct {
					Kind string `yaml:"kind"`
					Name string `yaml:"name"`
				} `yaml:"from"`
			} `yaml:"dockerStrategy"`
		} `yaml:"strategy"`
		Output struct {
			To struct {
				Kind string `yaml:"kind"`
				Name string `yaml:"name"`
			} `yaml:"to"`
		} `yaml:"output"`
		PostCommit struct {
			Script string `yaml:"script"`
		} `yaml:"postCommit"`
	} `yaml:"spec"`
}

// Kustomization struct
type Kustomization struct {
	APIVersion      string   `yaml:"apiVersion"`
	Kind            string   `yaml:"kind"`
	Bases           []string `yaml:"bases"`
	Namespace       string   `yaml:"namespace"`
	PatchesJSON6902 []struct {
		Path   string `yaml:"path"`
		Target struct {
			Group   string `yaml:"group"`
			Kind    string `yaml:"kind"`
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
		} `yaml:"target"`
	} `yaml:"patchesJson6902"`
}
