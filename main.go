package main

import (
	"fmt"
	"log"
	"net/http"

	"go_todos/internal/app/db"
	"go_todos/internal/app/tmp"
    "go_todos/internal/app/handler"


	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    status := db.NewRedisConn()
    if status != nil {
        log.Fatalf("Could not connect to Redis: %v", status)
    } else {
        fmt.Println("Connected to Redis successfully")
    }

	defer db.CloseDBConn()

	err := tmp.ParseTemplates()
	if err != nil {
		log.Panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	r.Get("/", handler.HandleGetTasks)
    r.Post("/tasks", handler.HandleCreateTask)
    r.Put("/tasks/{id}/toggle", handler.HandleToggleTask)
    r.Delete("/tasks/{id}", handler.HandleDeleteTask)
	r.Get("/tasks/{id}/edit", handler.HandleEditTask)
    r.Put("/tasks/{id}", handler.HandleUpdateTaks)

    fmt.Println("STARTING SERVER")

    http.ListenAndServe(":3242", r)
}
