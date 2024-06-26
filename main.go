package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/yash91989201/rss_aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT value is required.")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("Database URL is required.")
	}

	conn, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal("Unable to connect to database.", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           3000,
			},
		),
	)

	v1Router := chi.NewRouter()

	// endpoints
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Get("/users", apiCfg.handlerGetUser)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Starting server on port:%v", portString)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
