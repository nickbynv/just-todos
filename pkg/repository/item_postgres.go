package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "just-todos"
	"strings"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemPostgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db}
}

func (r *ItemPostgres) Create(listId int, item todo.Item) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf(
		`INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id`,
		itemsTable,
	)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)

	var itemId int
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf(
		`INSERT INTO %s (list_id, item_id) VALUES ($1, $2)`,
		listsItemsTable,
	)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return itemId, nil
}

func (r *ItemPostgres) GetAll(userId, listId int) ([]todo.Item, error) {
	var items []todo.Item

	query := fmt.Sprintf(
		`SELECT
    		i.id,
    		i.title,
    		i.description,
    		i.done
		FROM %s i
		JOIN %s li
			ON li.item_id = i.id
		JOIN %s ul
			ON ul.list_id = li.list_id
		WHERE li.list_id = $1
			AND ul.user_id = $2`,
		itemsTable,
		listsItemsTable,
		usersListsTable,
	)

	err := r.db.Select(&items, query, listId, userId)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemPostgres) GetById(userId, itemId int) (todo.Item, error) {
	var item todo.Item

	query := fmt.Sprintf(
		`SELECT
    		i.id,
    		i.title,
    		i.description,
    		i.done
		FROM %s i
		JOIN %s li
			ON li.item_id = i.id
		JOIN %s ul
			ON ul.list_id = li.list_id
		WHERE i.id = $1
			AND ul.user_id = $2`,
		itemsTable,
		listsItemsTable,
		usersListsTable,
	)

	err := r.db.Get(&item, query, itemId, userId)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (r *ItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s i
		USING
			%s li,
		    %s ul
		WHERE i.id = li.item_id
			AND li.list_id = ul.list_id
			AND ul.user_id = $1
			AND i.id = $2`,
		itemsTable,
		listsItemsTable,
		usersListsTable,
	)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *ItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	// title = $1
	// description = $1
	// title = $1, description = $2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s i
			SET %s
		FROM
		    %s li,
		    %s ul
		WHERE i.id = li.item_id
			AND li.list_id = ul.list_id
			AND ul.user_id = $%d
			AND i.id = $%d`,
		itemsTable,
		setQuery,
		listsItemsTable,
		usersListsTable,
		argId,
		argId+1,
	)

	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)

	return err
}
