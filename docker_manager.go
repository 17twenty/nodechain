package nodechain

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

type DockerManager struct {
	Image     string
	Workspace string // local directory to mount
	Container string // name or id assigned after start
	mu        sync.Mutex
	running   bool
}

func NewDockerManager(image, workspace string) *DockerManager {
	return &DockerManager{
		Image:     image,
		Workspace: workspace,
	}
}

func (m *DockerManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil
	}

	// Always attempt to remove existing container with this name
	exec.Command("docker", "rm", "-f", "nodechain-agent").Run()

	// Now run a fresh container
	cmd := exec.Command(
		"docker", "run",
		"-d",
		"--name=nodechain-agent",
		"-v", fmt.Sprintf("%s:/workspace", m.Workspace),
		m.Image,
		"sleep", "infinity",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %s (%v)", out, err)
	}

	m.Container = "nodechain-agent"
	m.running = true
	return nil
}

func (m *DockerManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	exec.Command("docker", "rm", "-f", m.Container).Run()
	m.running = false
	return nil
}

func (m *DockerManager) Run(cmdStr string) (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return "", "", fmt.Errorf("docker container not running")
	}

	args := append([]string{"exec", m.Container, "bash", "-lc"}, cmdStr)
	cmd := exec.Command("docker", args...)

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	return out.String(), errOut.String(), err
}
