package main

import (
	"ai-agent-business/internal/app"
	"ai-agent-business/internal/infra/logger"
	"ai-agent-business/internal/infra/variables"
)

func main() {
	logger.Init(&logger.Option{
		ServiceName:    variables.ServiceName(),
		ServiceVersion: variables.ServiceVersion(),
		Environment:    variables.Environment(),
		LogLevel:       variables.LogLevel(),
	})

	defer logger.Sync()

	application := app.Instance()
	application.Start(false)
}
