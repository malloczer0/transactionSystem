package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
	"transactionSystemTestTask/config"
	"transactionSystemTestTask/controller"
	libError "transactionSystemTestTask/lib/error"
	"transactionSystemTestTask/lib/validator"
	"transactionSystemTestTask/logger"
	"transactionSystemTestTask/pool"
	"transactionSystemTestTask/service"
	"transactionSystemTestTask/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	ctx := context.Background()
	cfg := config.Get()
	log := logger.Get()

	// Init repository storeObject (with postgresql inside)
	storeObject, err := store.New(ctx)
	if err != nil {
		return errors.Wrap(err, "storeObject.New failed")
	}

	// Init service manager
	serviceManager, err := service.NewManager(ctx, storeObject)
	if err != nil {
		return errors.Wrap(err, "manager.New failed")
	}

	// Init controllers
	transactionController := controller.NewTransactions(ctx, serviceManager, log)

	// Initialize Echo instance
	echoObject := echo.New()
	echoObject.Validator = validator.NewValidator()
	echoObject.HTTPErrorHandler = libError.Error
	// Disable Echo JSON log in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := echoObject.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	echoObject.Use(middleware.Logger())
	echoObject.Use(middleware.Recover())

	// API V1
	v1 := echoObject.Group("/v1")

	// Transaction routes
	transactionRoutes := v1.Group("/transactions")
	transactionRoutes.POST("/create", transactionController.Create)

	// Init client pool
	clientPool, _ := pool.Initialisation(ctx, storeObject)

	// Init pool handler
	go pool.Handler(ctx, *serviceManager, clientPool)

	// Start server
	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	echoObject.Logger.Fatal(echoObject.StartServer(server))

	return nil
}
