package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"gotest/app/controllers"
	"gotest/app"
)

func NewRouter(a *app.App) *mux.Router {
	r := mux.NewRouter()
	uc := controllers.New()
	api := r.PathPrefix("/api/v1").Subrouter()

	// Uploads
	// api.HandleFunc("/images/upload", middlewares.Logger(middlewares.RequireAuthentication(a, uploadController.UploadImage, true))).Methods(http.MethodPost)

	// Users
	api.HandleFunc("/users", uc.HelloUser).Methods(http.MethodGet)
	api.HandleFunc("/users", uc.Create).Methods(http.MethodPost)
	
	return r
}