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
	logger := logger.Get()

	// Init repository store object (with postgresql inside)
	log.Println("Init repository store object")
	storeObject, err := store.New(ctx)
	if err != nil {
		return errors.Wrap(err, "storeObject.New failed")
	}

	// Init service manager
	log.Println("Init service manager")
	serviceManager, err := service.NewManager(ctx, storeObject)
	if err != nil {
		return errors.Wrap(err, "manager.New failed")
	}

	// Init controllers
	log.Println("Init controllers")
	transactionController := controller.NewTransactions(ctx, serviceManager, logger)

	// Initialize Echo instance
	log.Println("Initialize Echo instance")
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
	log.Println("Init client pool")
	err = func() error {
		err = pool.Initialisation(ctx, storeObject)
		if err != nil {
			return errors.Wrap(err, "pool.Initialisation failed")
		}
		return nil
	}()

	// Init pool handler
	log.Println("Init pool handler")
	go pool.Handler(ctx, storeObject)
	if err != nil {
		return errors.Wrap(err, "pool.Handler failed")
	}

	// Start server
	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	echoObject.Logger.Fatal(echoObject.StartServer(server))

	return nil
}
