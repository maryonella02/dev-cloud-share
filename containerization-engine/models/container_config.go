package models

type ContainerConfig struct {
	Image       string            `json:"image"`
	Command     []string          `json:"command"`
	Environment map[string]string `json:"environment"`
	Memory      int64             `json:"memory"`
	NanoCPUs    int64             `json:"nano-cpus"`
}
