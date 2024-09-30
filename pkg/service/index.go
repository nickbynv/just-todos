package service

import (
	todo "just-todos"

	"just-todos/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type List interface {
	Create(userId int, list todo.List) (int, error)
	GetAll(userId int) ([]todo.List, error)
	GetById(userId, listId int) (todo.List, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todo.UpdateListInput) error
}

//type Item interface {
//
//}

type Service struct {
	Authorization
	List
	//Item
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		List:          NewListService(repos.List),
	}
}
