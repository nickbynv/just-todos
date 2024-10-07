package repository

import (
	todo "just-todos"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type List interface {
	Create(userId int, list todo.List) (int, error)
	GetAll(userId int) ([]todo.List, error)
	GetById(userId, listId int) (todo.List, error)
	Delete(userId, listId int) error
	Update(userId, listId int, list todo.UpdateListInput) (todo.List, error)
}

type Item interface {
	Create(listId int, item todo.Item) (int, error)
	GetAll(userId, listId int) ([]todo.Item, error)
	GetById(userId, itemId int) (todo.Item, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, list todo.UpdateItemInput) error
}

type Repository struct {
	Auth
	List
	Item
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth: NewAuthPostgres(db),
		List: NewListPostgres(db),
		Item: NewItemPostgres(db),
	}
}
