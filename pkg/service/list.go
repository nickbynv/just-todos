package service

import (
	todo "just-todos"
	"just-todos/pkg/repository"
)

type ListService struct {
	repo repository.List
}

func (s *ListService) Create(userId int, list todo.List) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *ListService) GetAll(userId int) ([]todo.List, error) {
	return s.repo.GetAll(userId)
}

func (s *ListService) GetById(userId, listId int) (todo.List, error) {
	return s.repo.GetById(userId, listId)
}

func (s *ListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *ListService) Update(userId, listId int, input todo.UpdateListInput) (todo.List, error) {
	if err := input.Validate(); err != nil {
		return todo.List{}, err
	}

	return s.repo.Update(userId, listId, input)
}

func NewListService(repo repository.List) *ListService {
	return &ListService{repo}
}
