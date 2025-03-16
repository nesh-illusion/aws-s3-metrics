package routes

import (
	"project_security_one/api/handlers"

	"github.com/gorilla/mux"
)

// SetupRouter initializes API routes
func SetupRouter(userHandler *handlers.UserHandler) *mux.Router {
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/users/register", userHandler.RegisterUser).Methods("POST")

	return router
}
