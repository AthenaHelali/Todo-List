package task

import (
	"fmt"
	"todo-list/models"
)

type ServiceRepository interface {
	// CreateNewTask DoesThisUserHasThisCategoryID(userID, categoryID int) (bool, error)
	CreateNewTask(t models.Task) (models.Task, error)
	ListUserTask(userID int) ([]models.Task, error)
}

type Service struct {
	repository ServiceRepository
}

func NewService(rep ServiceRepository) *Service {
	return &Service{
		repository: rep,
	}
}

type CreateTaskRequest struct {
	Title               string
	DueDate             string
	CategoryID          int
	AuthenticatedUserID int
}
type CreateResponse struct {
	Task models.Task
}

func (s Service) CreateTask(req CreateTaskRequest) (CreateResponse, error) {

	/*if ok, _ := s.repository.DoesThisUserHasThisCategoryID(req.AuthenticatedUserID, req.CategoryID); !ok {
		return CreateResponse{}, fmt.Errorf("category-id %d is not found", req.CategoryID)
	}*/
	createdTask, cErr := s.repository.CreateNewTask(models.Task{
		Title:      req.Title,
		DueDate:    req.DueDate,
		CategoryID: req.CategoryID,
		IsDone:     false,
		UserID:     req.AuthenticatedUserID,
	})
	if cErr != nil {
		return CreateResponse{}, fmt.Errorf("can't create new task: %v", cErr)
	}

	return CreateResponse{createdTask}, nil

}

type ListRequest struct {
	UserID int
}
type ListResponse struct {
	Tasks []models.Task
}

func (s Service) ListTask(req ListRequest) (ListResponse, error) {
	tasks, err := s.repository.ListUserTask(req.UserID)
	if err != nil {
		return ListResponse{}, fmt.Errorf("can't list tasks: %v", err)
	}

	return ListResponse{Tasks: tasks}, nil
}
