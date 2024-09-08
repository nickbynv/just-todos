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

//type List interface {
//
//}

//type Item interface {
//
//}

type Service struct {
	Authorization
	//List
	//Item
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
