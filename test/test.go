package test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"ai-agent-business/internal/app"
	"ai-agent-business/internal/infra/logger"
	"ai-agent-business/internal/infra/variables"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

func InitLogger() {
	logger.Init(&logger.Option{
		ServiceName:    variables.ServiceName(),
		ServiceVersion: variables.ServiceVersion(),
		Environment:    variables.Environment(),
		LogLevel:       variables.LogLevel(),
	})
}

func ComposeUp(t *testing.T) {
	result := make(chan bool, 1)
	defer close(result)

	go func() {
		Shell(t, fmt.Sprintf("cd %s && make compose-up", findProjectFolder(t)))

		for {
			if containerHealthy("postgres") && containerHealthy("redis") {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}

		result <- true
	}()

	select {
	case <-result:
	case <-time.After(10 * time.Second):
	}
}

func ComposeDown(t *testing.T) {
	Shell(t, fmt.Sprintf("cd %s && make compose-up", findProjectFolder(t)))
}

func StartApplication() {
	if app.Instance().IsRunning() {
		return
	}

	app.Instance().Start(true)
}

func StopApplication() {
	app.Instance().Stop()
	_, _ = ShellErr("fuser -k 5000/tcp")
}

func Request() *baloo.Client {
	return baloo.New("http://localhost:5000")
}

func Shell(t *testing.T, command string) string {
	out, err := ShellErr(command)

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			assert.FailNow(t, fmt.Sprintf("out: %s\nerr: %s", out, string(ee.Stderr)))
		}

		assert.FailNow(t, err.Error())
	}

	return out
}

func ShellErr(command string) (string, error) {
	out, err := exec.Command("/bin/sh", "-c", command).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func findProjectFolder(t *testing.T) string {
	folders := strings.Split(Shell(t, "pwd"), "/")
	for i := len(folders) - 1; i >= 0; i-- {
		folder := strings.Join(folders[:i], "/")

		out, err := ShellErr(fmt.Sprintf("ls %s | grep go.mod", folder))
		if err != nil {
			continue
		}

		if strings.TrimSpace(out) == "go.mod" {
			return folder
		}
	}

	assert.FailNow(t, "Project folder not found")
	return ""
}

func containerHealthy(image string) bool {
	container, _ := ShellErr(fmt.Sprintf("docker ps | grep %s", image))
	return strings.Contains(container, "healthy")
}
