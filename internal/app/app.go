package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"ai-agent-business/internal/app/adapters"
	"ai-agent-business/internal/app/domain"
	"ai-agent-business/internal/infra/aws"
	"ai-agent-business/internal/infra/database"
	"ai-agent-business/internal/infra/logger"
	"ai-agent-business/internal/infra/logger/attributes"
	"ai-agent-business/internal/infra/server"
	"ai-agent-business/internal/infra/variables"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
)

type App struct {
	running bool
	locker  sync.Mutex

	server    *echo.Echo
	handlers  *adapters.Handlers
	services  *domain.Services
	databases *database.Databases
}

var app = new(App)

func Instance() *App {
	return app
}

func (app *App) Start(async bool) {
	app.locker.Lock()

	if app.running {
		app.locker.Unlock()
		return
	}

	start := time.Now()
	logger.Info(context.Background(), fmt.Sprintf("Starting application %s:%s", variables.ServiceName(), variables.ServiceVersion()), nil)

	app.build()

	if async {
		go app.startServer(start)
	} else {
		app.startServer(start)
	}
}

func (app *App) Stop() {
	app.locker.Lock()

	if !app.running {
		app.locker.Unlock()
		return
	}

	defer app.setRunning(false)
	defer app.locker.Unlock()

	logger.Warn(context.Background(), "Stopping application", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "Error while trying to close echo server", attributes.New().WithError(err))
	}

	app.databases.Close()
	app.dispose()

	logger.Warn(context.Background(), "Application stopped", nil)
}

func (app *App) IsRunning() bool {
	return app.running
}

func (app *App) startServer(start time.Time) {
	defer app.setRunning(false)
	go func() {
		app.printElapsed(start)
		app.locker.Unlock()
	}()

	if variables.IsLambda() {
		lambdaAdapter := &aws.LambdaAdapter{Echo: app.server}
		lambda.Start(lambdaAdapter.Handler)
		logger.Warn(context.Background(), "Application stopped [Lambda]", nil)
	} else {
		// Set up graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			app.Stop()
		}()

		if err := app.server.Start(fmt.Sprintf("%s:%d", variables.ServerHost(), variables.ServerPort())); err != nil {
			logger.Warn(context.Background(), "Application stopped", attributes.New().WithError(err))
		}
	}
}

func (app *App) build() {
	app.databases = database.NewDatabases()
	app.services = domain.NewServices(app.databases)
	app.handlers = adapters.NewHandlers(app.services)
	app.server = server.New()
	app.handlers.Configure(app.server)
}

func (app *App) dispose() {
	app.server = nil
	app.handlers = nil
	app.services = nil
	app.databases = nil
}

func (app *App) printElapsed(start time.Time) {
	elapsed := time.Since(start)
	logger.Info(context.Background(), fmt.Sprintf("Application %s:%s started in %v", variables.ServiceName(), variables.ServiceVersion(), elapsed.String()), nil)
}

func (app *App) setRunning(run bool) {
	app.running = run
}
