package handler

import (
	"log"
	"net/http"
	"strconv"

	"go_todos/internal/app/db"
	"go_todos/internal/app/models"
	"go_todos/internal/app/tmp"

	"github.com/go-chi/chi/v5"
)

func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	items, err := db.FetchTasks()

	if err != nil {
		log.Printf("Error fetching tasks %v", err)

		return
	}

	count := len(items)

	log.Printf("Tasks found %d", count)

	var competed int = 0
	for _, item := range items {
		if item.Status == true {
			competed += 1
		}
	}

	log.Printf("Completed tasks found %d", competed)

	data := models.TaskList{
		Tasks: items,
		Count: count,
		CompletedCount: competed,
	}

	tmp.TMPL.ExecuteTemplate(w, "Base", data)
}


func HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")

	if title == "" {
		tmp.TMPL.ExecuteTemplate(w, "Form", nil)
		return
	}

	newTask, err := db.InsertTask(title)
	if err != nil {
		log.Printf("Error inserting taks %v", err)
	}

	items, err := db.FetchTasks()

	if err != nil {
		log.Printf("Error fetching tasks %v", err)

		return
	}

	count := len(items)

	log.Printf("Tasks found %d", count)

	w.WriteHeader(http.StatusCreated)
	tmp.TMPL.ExecuteTemplate(w, "Form", nil)

	tmp.TMPL.ExecuteTemplate(w, "TotalCount", map[string]any{"Count": count, "SwapOOB": true})

	tmp.TMPL.ExecuteTemplate(w, "Item", map[string]any{"Item": newTask, "SwapOOB": true})
}

func HandleToggleTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		log.Printf("Error parsing id into int %w", err)
		return
	}

	err = db.ToggleTask(id)

	if err != nil {
		log.Printf("Error toggling task %w", err)
		return
	}

	items, err := db.FetchTasks()

	if err != nil {
		log.Printf("Error fetching tasks %v", err)

		return
	}

	var competed int = 0
	for _, item := range items {
		if item.Status == true {
			competed += 1
		}
	}

	tmp.TMPL.ExecuteTemplate(w, "CompletedCount", map[string]any{"Count": competed, "SwapOOB": true})
}


func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		log.Printf("Error parsing id into int %w", err)
		return
	}

	err = db.RemoveTask(id)

	if err != nil {
		log.Printf("Error fetching tasks %v", err)

		return
	}

	items, err := db.FetchTasks()

	if err != nil {
		log.Printf("Error fetching tasks %v", err)

		return
	}

	count := len(items)

	var competed int = 0
	for _, item := range items {
		if item.Status == true {
			competed += 1
		}
	}

	tmp.TMPL.ExecuteTemplate(w, "TotalCount", map[string]any{"Count": count, "SwapOOB": true})
	tmp.TMPL.ExecuteTemplate(w, "CompletedCount", map[string]any{"Count": competed, "SwapOOB": true})

	w.WriteHeader(http.StatusOK)
}


func HandleEditTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		log.Printf("Error parsing id into int %w", err)
		return
	}

	task, err := db.FetchTask(id)

	if err != nil {
		log.Printf("Error fetching tasks %v", err)

		return
	}

	tmp.TMPL.ExecuteTemplate(w, "Item", map[string]any{"Item": task, "Editing": true})
}

func HandleUpdateTaks(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		log.Printf("Error parsing id into int %w", err)
		return
	}

	newTitle := r.FormValue("title")

	if newTitle == "" {
		task, err := db.FetchTask(id)

		if err != nil {
			log.Printf("Error fetching tasks %v", err)
	
			return
		}

		tmp.TMPL.ExecuteTemplate(w, "Item", map[string]any{"Item": task, "Editing": false})

		return	
	}

	task, err := db.UpdateTask(id, newTitle)
	if err != nil {
		log.Printf("Error updating tasks %v", err)
		return
	}

	tmp.TMPL.ExecuteTemplate(w, "Item", map[string]any{"Item": task, "Editing": false})
}