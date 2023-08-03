package memorystore

import (
	"todo-list/models"
)

type Task struct {
	tasks []models.Task
}

func NewTaskStore() *Task {
	return &Task{
		make([]models.Task, 0),
	}
}

func (t *Task) CreateNewTask(task models.Task) (models.Task, error) {
	task.ID = len(t.tasks) + 1
	t.tasks = append(t.tasks, task)
	return task, nil

}

func (t *Task) ListUserTask(userID int) ([]models.Task, error) {
	var userTasks []models.Task
	for _, task := range t.tasks {
		if task.UserID == userID {
			userTasks = append(userTasks, task)
		}
	}
	return userTasks, nil
}
