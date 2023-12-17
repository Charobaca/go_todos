package db

import (
    "context"
    "fmt"
    "encoding/json"

    "go_todos/internal/app/models"

    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var DB *redis.Client

func NewRedisConn() error {
    db := redis.NewClient(&redis.Options{
        Addr:     "redis:6379",
        DB:       1,
    })

    _, err := db.Ping(ctx).Result()
    if err != nil {
        return err
    }

    DB = db

    return nil
}

func CloseDBConn() error {
    err := DB.Close()

    if err != nil {
        return err
    }

    fmt.Println("Redis connection closed successfully")

    return nil
}


func FetchTasks() ([]models.Task, error) {
    var tasks []models.Task

    keyExists, err := DB.Exists(ctx, "tasks").Result()
    if err != nil {
        return nil, err
    }
    if keyExists == 0 {
        // If the "tasks" key doesn't exist, initialize it with an empty tasks slice
        jsonData, err := json.Marshal(tasks)
        if err != nil {
            return nil, err
        }

        err = DB.Set(ctx, "tasks", jsonData, 0).Err()
        if err != nil {
            return nil, err
        }
        return tasks, nil
    }

    jsonData, err := DB.Get(ctx, "tasks").Result()
    if err != nil {
        fmt.Println("Get failed", err)
        return nil, err
    }

    err = json.Unmarshal([]byte(jsonData), &tasks)
    if err != nil {
        fmt.Println("Unmarshal failed", err)
		return nil, err
    }

    return tasks, nil
}


func FetchTask(ID int) (models.Task, error) {
    var tasks []models.Task

    tasks, _ = FetchTasks()

    if tasks == nil {
        return models.Task{}, nil
    }

    var targetTask models.Task 

    for _, task := range tasks {
        if task.ID == ID {
            targetTask = models.Task(task)
        }
    }

    return targetTask, nil
}

func InsertTask(title string) (models.Task, error) {
	// Fetch existing tasks
	tasks, err := FetchTasks()
	if err != nil {
		return models.Task{}, err
	}

	// Generate a new unique ID
	newID := generateUniqueID(tasks)

	// Create a new task with the provided title and the generated ID
	newTask := models.Task{
		ID:    newID,
		Title: title,
	}

	// Append the new task to the existing tasks
	tasks = append(tasks, newTask)

	// Update the tasks in the database
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return models.Task{}, err
	}

	err = DB.Set(ctx, "tasks", jsonData, 0).Err()
	if err != nil {
		return models.Task{}, err
	}

	return newTask, nil
}

func generateUniqueID(tasks []models.Task) int {
	// Find the maximum ID currently in use
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	// Increment the maximum ID to generate a new unique ID
	return maxID + 1
}



func RemoveTask(ID int) error {
	tasks, err := FetchTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	err = DB.Set(ctx, "tasks", jsonData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func UpdateTask(ID int, newTitle string) (models.Task, error) {
	// Fetch existing tasks
	tasks, err := FetchTasks()
	if err != nil {
		return models.Task{}, err
	}

	// Find the task with the specified ID
	found := false
    foundedIndex := 0
	for i, task := range tasks {
		if task.ID == ID {
			// Update the title of the found task
			tasks[i].Title = newTitle
			found = true
            foundedIndex = i
			break
		}
	}

	// If the task with the specified ID is not found, return an error
	if !found {
		fmt.Errorf("Task with ID %d not found", ID)
        return models.Task{}, err
	}

	// Update the tasks in the database
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return models.Task{}, err
	}

	err = DB.Set(ctx, "tasks", jsonData, 0).Err()
	if err != nil {
		return models.Task{}, err
	}

	return tasks[foundedIndex], nil
}


func ToggleTask(ID int) error {
	// Fetch existing tasks
	tasks, err := FetchTasks()
	if err != nil {
		return err
	}

	// Find the task with the specified ID
	found := false
	for i, task := range tasks {
		if task.ID == ID {
			// Toggle the status of the found task
			tasks[i].Status = !tasks[i].Status
			found = true
			break
		}
	}

	// If the task with the specified ID is not found, return an error
	if !found {
		return fmt.Errorf("Task with ID %d not found", ID)
	}

	// Update the tasks in the database
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	err = DB.Set(ctx, "tasks", jsonData, 0).Err()
	if err != nil {
		return err
	}
 
	return nil
}