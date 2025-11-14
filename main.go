package main

import (
	"log"
	"net/http"

	"github.com/deside01/is_available/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	api := chi.NewRouter()
	v1 := chi.NewRouter()
	v1.Get("/check", handlers.Check)
	v1.Get("/data", handlers.GetData)

	api.Mount("/v1", v1)
	r.Mount("/api", api)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
		Handler:      r,
	}
	// addr := fmt.Sprintf("%v:%v", "localhost", "8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
