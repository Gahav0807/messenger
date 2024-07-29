package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
    // Параметры подключения к базе данных
    PGUSER := "postgres"
    PGPASSWORD := "vahagn2009"
    PGDATABASE := "messanger"
    PGHOST := "localhost"
    dbPort := 5432

    // Создание строки подключения
    connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    PGHOST, dbPort, PGUSER, PGPASSWORD, PGDATABASE)

    // Открытие подключения к базе данных
    db, err := sql.Open("postgres", connString)
    if err != nil {
        log.Fatalf("Ошибка подключения к базе данных: %v", err)
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        log.Fatalf("Ошибка проверки подключения: %v", err)
    }

    fmt.Println("Успешное подключение к базе данных!")

    // Здесь можно выполнять запросы к базе данных
    // ...

    rows, err := db.Query("SELECT password, created_at, username, id FROM Users")
    if err != nil {
        log.Fatalf("Ошибка выполнения запроса: %v", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var password, createdAt, username string
        var id int
        err := rows.Scan(&password, &createdAt, &username, &id)
        if err != nil {
            log.Fatalf("Ошибка чтения результатов: %v", err)
        }
        fmt.Printf("ID: %d, Username: %s, Password: %s, Created At: %s\n", id, username, password, createdAt)
    }
    
    err = rows.Err()
    if err != nil {
        log.Fatalf("Ошибка во время итерации результатов: %v", err)
    }
}
