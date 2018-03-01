package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"gotest/app/users"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	// uc := users.HelloUser()
	api := r.PathPrefix("/api/v1").Subrouter()

	// Uploads
	// api.HandleFunc("/images/upload", middlewares.Logger(middlewares.RequireAuthentication(a, uploadController.UploadImage, true))).Methods(http.MethodPost)

	// Users
	api.HandleFunc("/users", users.HelloUser).Methods(http.MethodGet)
	
	return r
}