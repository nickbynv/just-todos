package service

import (
	todo "just-todos"

	"just-todos/pkg/repository"
)

type Auth interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type List interface {
	Create(userId int, list todo.List) (int, error)
	GetAll(userId int) ([]todo.List, error)
	GetById(userId, listId int) (todo.List, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateListInput) (todo.List, error)
}

type Item interface {
	Create(userId, listId int, item todo.Item) (int, error)
	GetAll(userId, listId int) ([]todo.Item, error)
	GetById(userId, itemId int) (todo.Item, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}

type Service struct {
	Auth
	List
	Item
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repo.Auth),
		List: NewListService(repo.List),
		Item: NewItemService(repo.Item, repo.List),
	}
}
