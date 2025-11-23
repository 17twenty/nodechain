package nodechain

import "fmt"

type DockerExecTool struct {
	Docker *DockerManager
}

func (t *DockerExecTool) Name() string { return "docker_exec" }

func (t *DockerExecTool) Run(input any) (any, error) {
	cmdStr, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("docker_exec: input must be string")
	}

	stdout, stderr, err := t.Docker.Run(cmdStr)
	if err != nil {
		return map[string]any{
			"stdout": stdout,
			"stderr": stderr,
			"error":  err.Error(),
		}, nil
	}

	return map[string]any{
		"stdout": stdout,
		"stderr": stderr,
	}, nil
}
