package main

import (
  "os"
  "fmt"
	"net/http"
	"github.com/gorilla/handlers"
  log "github.com/sirupsen/logrus"
  "gotest/app/routes"
  "gotest/config"
  "gotest/database"
)

type App struct {
	Config   config.Config
	Database *database.MySQLDB
	Redis    *database.RedisDB
}

func init() {
  // Log as JSON instead of the default ASCII formatter.
  log.SetFormatter(&log.JSONFormatter{})

  // Output to stdout instead of the default stderr
  // Can be any io.Writer, see below for File example
  log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  log.SetLevel(log.WarnLevel)
}

func main() {
	cfg, err := config.New("config/config.json")
	if (err != nil) {
		log.Fatal(err)
	}

	db, err := database.NewMySQLDB(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}

	r := routes.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	port := 3002
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("APP is listening on port: %d\n", port)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
