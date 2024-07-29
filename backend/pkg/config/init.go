package config

import (
	"fmt"
	"os"

	"my-go-project/pkg/core"
    "my-go-project/internal/handlers"
    "go.uber.org/zap"
)

func InitDependencies() (*zap.Logger, *database.Database) {
    // Initialize logger
    logger := initLogger()

    // Set logger in other parts of the program
    database.SetLogger(logger)
    handlers.SetLogger(logger)

    // Initialize database
    database := database.NewDatabase()

    return logger, database
}

func ShutdownDependencies(logger *zap.Logger, database *database.Database) {
    // Закрытие пула с бд, сохранение логов при завершении программы(ctrl+c)
    database.CloseConnection()
    logger.Sync()
}

func initLogger() *zap.Logger {
    logger, err := zap.NewProduction()
    if err != nil {
        fmt.Printf("Can't initialize zap logger: %v", err)
        os.Exit(1)
    }
    return logger
}
