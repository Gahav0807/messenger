// routers
package api

import (
	"fmt"
	"net/http"
    "my-go-project/internal/handlers"
    "github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
    router := mux.NewRouter()

    authHandler := handlers.NewAuthHandler()

    router.HandleFunc("/", handleRoot).Methods("GET")
    // auth,register
    router.HandleFunc("/tryAuth/{username}/{password}", authHandler.Login).Methods("GET")
    router.HandleFunc("/tryRegister/{username}/{password}", authHandler.Registration).Methods("GET")
    
    return router
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Println("New request received!")
    w.Write([]byte("Welcome to the clicker-webapp server!"))
}
