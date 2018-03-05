package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gotest/app/controllers"
	"gotest/app"
	"gotest/app/middlewares"
	"gotest/app/services"
	"gotest/app/models"
)

func NewRouter(a *app.App) *mux.Router {
	r := mux.NewRouter()
	uc := controllers.New(a.Database)
	api := r.PathPrefix("/api/v1").Subrouter()

	// Services
	jwtAuth := services.NewJWTAuthService(&a.Config.JWT, a.Redis)

	// Controller
	uh := models.NewUserHelper(a.Database)
	ac := controllers.NewAuthController(a, uh, jwtAuth)

	// Uploads
	// api.HandleFunc("/images/upload", middlewares.Logger(middlewares.RequireAuthentication(a, uploadController.UploadImage, true))).Methods(http.MethodPost)

	// graphql
	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    models.RootQuery,
		Mutation: nil,
	})
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: true,
	})
	r.Handle("/graphiql", h)
	r.Handle("/graphql", middlewares.RequireAuthenGraphql(a, h))

	// Users
	api.HandleFunc("/users", middlewares.Logger(uc.HelloUser)).Methods(http.MethodGet)
	api.HandleFunc("/users", middlewares.Logger(uc.Create)).Methods(http.MethodPost)

	// Authentication
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", middlewares.Logger(ac.Authenticate)).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", middlewares.Logger(middlewares.RequireRefreshToken(a, ac.RefreshTokens))).Methods(http.MethodGet)
	//auth.HandleFunc("/update", middlewares.Logger(middlewares.RequireAuthentication(a, uc.Update, false))).Methods(http.MethodPut)
	auth.HandleFunc("/logout", middlewares.Logger(middlewares.RequireAuthentication(a, ac.Logout, false))).Methods(http.MethodGet)
	auth.HandleFunc("/logout/all", middlewares.Logger(middlewares.RequireAuthentication(a, ac.LogoutAll, false))).Methods(http.MethodGet)

	auth.HandleFunc("/facebook/login", ac.FaceBookLogin).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/callback", ac.FaceBookCallback).Methods(http.MethodGet)
	auth.HandleFunc("/google/login", ac.GoogleLogin).Methods(http.MethodGet)
	auth.HandleFunc("/google/callback", ac.GoogleCallBack).Methods(http.MethodGet)
	return r
}