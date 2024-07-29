package auth

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/rs/cors"
    "go.uber.org/zap"

    "my-go-project/api"
    "my-go-project/pkg/config"
)

func Start() {
    // Initialize dependencies
    logger, database := config.InitDependencies()
    defer config.ShutdownDependencies(logger, database)

    // Setup router
    router := api.SetupRouter()

    // Настраиваем CORS
    c := cors.New(cors.Options{
        AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:     []string{"Content-Type", "Authorization"},
        AllowedOrigins:     []string{"*"},
        AllowCredentials:   true,
        OptionsPassthrough: false,
    })
    handler := c.Handler(router)

    // Закрытие пула с бд, сохранение логов при завершении программы(ctrl+c)
    shutdownChan := make(chan os.Signal, 1)
    signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-shutdownChan
        logger.Info("Shutting down server...")

        // Закрытие пула с бд, сохранение логов перед завершением программы
        database.CloseConnection()
        logger.Sync()

        os.Exit(0)
    }()

    // Start server
    addr := ":9000"
    fmt.Printf("Server is running on http://localhost%s\n", addr)
    if err := http.ListenAndServe(addr, handler); err != nil {
        logger.Fatal("Server failed to start", zap.Error(err))
    }
}
