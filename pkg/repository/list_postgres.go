package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "just-todos"
	"strings"
)

type ListPostgres struct {
	db *sqlx.DB
}

func NewListPostgres(db *sqlx.DB) *ListPostgres {
	return &ListPostgres{db}
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
		`SELECT
    		l.id,
    		l.title,
    		l.description
		FROM %s l
		INNER JOIN %s ul
			ON l.id = ul.list_id
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
		`SELECT
    		l.id,
    		l.title,
    		l.description
		FROM %s l
		INNER JOIN %s ul
			ON l.id = ul.list_id
		WHERE ul.user_id = $1
			AND ul.list_id = $2`,
		listsTable,
		usersListsTable,
	)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *ListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s l
		USING %s ul
		WHERE l.id = ul.list_id
			AND ul.user_id = $1
			AND ul.list_id = $2`,
		listsTable,
		usersListsTable,
	)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *ListPostgres) Update(userId, listId int, input todo.UpdateListInput) (todo.List, error) {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// title = $1
	// description = $1
	// title = $1, description = $2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s l
			SET %s
			FROM %s ul
			WHERE l.id = ul.list_id
				AND ul.list_id = $%d
				AND ul.user_id = $%d
		RETURNING l.id, title, description`,
		listsTable,
		setQuery,
		usersListsTable,
		argId,
		argId+1,
	)

	args = append(args, listId, userId)

	row := r.db.QueryRow(query, args...)

	var list todo.List
	if err := row.Scan(&list.Id, &list.Title, &list.Description); err != nil {
		return todo.List{}, err
	}

	return list, nil
}
