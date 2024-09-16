package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "just-todos"
)

type ListPostgres struct {
	db *sqlx.DB
}

func (r *ListPostgres) Create(userId int, list todo.List) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf(
		"INSERT INTO %s(title, description) VALUES($1, $2) RETURNING id",
		listsTable,
	)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	var id int
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf(
		"INSERT INTO %s(user_id, list_id) VALUES($1, $2)",
		usersListsTable,
	)
	_, err = tx.Exec(createUsersListsQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *ListPostgres) GetAll(userId int) ([]todo.List, error) {
	var lists []todo.List

	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description
				FROM %s tl
				INNER JOIN %s ul
					ON tl.id = ul.list_id
						WHERE ul.user_id = $1`,
		listsTable,
		usersListsTable,
	)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *ListPostgres) GetById(userId, listId int) (todo.List, error) {
	var list todo.List

	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description
				FROM %s tl
				INNER JOIN %s ul
					ON tl.id = ul.list_id
						WHERE ul.user_id = $1
							AND ul.list_id = $2`,
		listsTable,
		usersListsTable,
	)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func NewListPostgres(db *sqlx.DB) *ListPostgres {
	return &ListPostgres{db: db}
}
