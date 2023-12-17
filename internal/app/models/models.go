package models

type Task struct {
	ID  	int
	Title   string
	Status  bool
}

type TaskList struct {
	Tasks 		   []Task
	Count		   int
	CompletedCount int
}