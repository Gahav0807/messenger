package handlers

import (
	"fmt"
	"net/http"
	"time"

	"my-go-project/pkg/core"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetLogger(l *zap.Logger) {
    logger = l
}
type AuthHandler struct {
	data *database.Database
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		data: database.NewDatabase(),
	}
}

// Метод для регистрации пользователя
// Получаем данные передаваемые в запросе, проверяем:
// Есть ли пользователь с таким же именем.Если есть - просим пользователя изменить ник
// Нету - записываем в бд и оповещаем об успехе
func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userName := vars["username"]
	password := vars["password"]

	logger.Info("Получаем инфу про пользователя из бд")

	query := fmt.Sprintf("SELECT * FROM Users WHERE username = '%s'", userName)
	result, err := h.data.GetData(ctx, query)
	if err != nil {
		logger.Error("Ошибка при получении информации о пользователе", zap.Error(err))
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	if len(result) > 0 {
		logger.Info("Пользователь с таким именем уже есть, оповещаем об этом")
		http.Error(w, "Invalid request method", http.StatusBadRequest)
        return
	}

    currentTime := time.Now()
    createdAt := currentTime.Format("02.01.2006")

	logger.Info("Пользователя с таким именем еще нет, заносим в бд")
	insertQuery := fmt.Sprintf("INSERT INTO Users (username, password, created_at) VALUES ('%s', '%s', '%s')", userName, password, createdAt)
	err = h.data.SetData(ctx, insertQuery)
	if err != nil {
		logger.Error("Ошибка при добавлении нового пользователя в бд", zap.Error(err))
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    vars := mux.Vars(r)
    userName := vars["username"]
    password := vars["password"]

    logger.Info("Получаем информацию о пользователе из базы данных")

    query := fmt.Sprintf("SELECT * FROM Users WHERE username = '%s'", userName)
    result, err := h.data.GetData(ctx, query)
    if err != nil {
        logger.Error("Ошибка при получении информации о пользователе", zap.Error(err))
        http.Error(w, "Database query error", http.StatusInternalServerError)
        return
    }

    if len(result) == 0 {
        logger.Info("Пользователь с таким именем не найден")
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    user := result[0]
    if user["password"] != password {
        logger.Info("Неверный пароль")
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    logger.Info("Пользователь успешно аутентифицирован")

    // Здесь можно сгенерировать и отправить токен для аутентификации пользователя
    w.WriteHeader(http.StatusOK)
}
