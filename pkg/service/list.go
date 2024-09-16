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

func NewListService(repo repository.List) *ListService {
	return &ListService{repo}
}
