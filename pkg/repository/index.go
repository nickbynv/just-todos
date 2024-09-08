package repository

import (
	todo "just-todos"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

//type List interface {
//
//}

//type Item interface {
//
//}

type Repository struct {
	Authorization
	//List
	//Item
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
