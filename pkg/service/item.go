package service

import (
	todo "just-todos"
	"just-todos/pkg/repository"
)

type ItemService struct {
	repo     repository.Item
	listRepo repository.List
}

func NewItemService(repo repository.Item, listRepo repository.List) *ItemService {
	return &ItemService{repo, listRepo}
}

func (s *ItemService) Create(userId, listId int, item todo.Item) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *ItemService) GetAll(userId, listId int) ([]todo.Item, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *ItemService) GetById(userId, itemId int) (todo.Item, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *ItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *ItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
