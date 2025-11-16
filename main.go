package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deside01/is_available/internal/config"
	"github.com/deside01/is_available/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
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
		Addr:         config.Data.Address,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  0,
		Handler:      r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running: %v", config.Data.Address)
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()
	<-stop

	log.Println("Получен сигнал остановки сервера")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}
