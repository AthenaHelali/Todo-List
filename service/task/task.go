package service

import (
	"fmt"
	"todo-list/models"
)

type TaskRepository interface {
	DoesThisUserHasThisCategoryID(userID, categoryID int) (bool, error)
	CreateNewTask(t models.Task) (models.Task, error)
}

type Task struct {
	repository TaskRepository
}
type CreateTaskRequest struct {
	Title               string
	DueDate             string
	CategoryID          int
	AuthenticatedUserID int
}
type CreateTaskResponse struct {
	Task models.Task
}

func (t Task) CreateTask(req CreateTaskRequest) (CreateTaskResponse, error) {

	if ok, _ := t.repository.DoesThisUserHasThisCategoryID(req.AuthenticatedUserID, req.CategoryID); !ok {
		return CreateTaskResponse{}, fmt.Errorf("category-id %d is not found", req.CategoryID)
	}
	createdTask, cErr := t.repository.CreateNewTask(models.Task{
		Title:      req.Title,
		DueDate:    req.DueDate,
		CategoryID: req.CategoryID,
		IsDone:     false,
		UserID:     req.AuthenticatedUserID,
	})
	if cErr != nil {
		return CreateTaskResponse{}, fmt.Errorf("can't create new task: %v", cErr)
	}

	return CreateTaskResponse{createdTask}, nil

}
func () 
