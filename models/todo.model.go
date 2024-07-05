package models

import (
	"errors"
	"fmt"
	"time"
)

type Todo struct {
	ID          uint64    `json:"id"`
	CreatedBy   uint64    `json:"created_by"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      bool      `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

func (t *Todo) GetAllTodos() ([]Todo, error) {
	query := fmt.Sprintf("SELECT id, title, description, status FROM todos WHERE created_by = %d ORDER BY created_at DESC", t.CreatedBy)

	rows, err := db.Query(query)
	if err != nil {
		return []Todo{}, err
	}
	// We close the resource
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)

		todos = append(todos, *t)
	}

	return todos, nil
}

func (t *Todo) GetNoteById() (Todo, error) {

	query := `SELECT id, title, description, status, created_at FROM todos
		WHERE created_by = ? AND id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return Todo{}, err
	}

	defer stmt.Close()

	var recoveredTodo Todo
	err = stmt.QueryRow(
		t.CreatedBy, t.ID,
	).Scan(
		&recoveredTodo.ID,
		&recoveredTodo.Title,
		&recoveredTodo.Description,
		&recoveredTodo.Status,
		&recoveredTodo.CreatedAt,
	)
	if err != nil {
		return Todo{}, err
	}

	return recoveredTodo, nil
}

func (t *Todo) CreateTodo() (Todo, error) {

	query := `INSERT INTO todos (created_by, title, description)
		VALUES(?, ?, ?) RETURNING *`

	stmt, err := db.Prepare(query)
	if err != nil {
		return Todo{}, err
	}

	defer stmt.Close()

	var newTodo Todo
	err = stmt.QueryRow(
		t.CreatedBy,
		t.Title,
		t.Description,
	).Scan(
		&newTodo.ID,
		&newTodo.CreatedBy,
		&newTodo.Title,
		&newTodo.Description,
		&newTodo.Status,
		&newTodo.CreatedAt,
	)
	if err != nil {
		return Todo{}, err
	}

	/* if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("error: an affected row was expected")
	} */

	return newTodo, nil
}
func (t *Todo) UpdateTodo() (Todo, error) {

	query := `UPDATE todos SET title = ?,  description = ?, status = ?
		WHERE created_by = ? AND id=? RETURNING id, title, description, status`

	stmt, err := db.Prepare(query)
	if err != nil {
		return Todo{}, err
	}

	defer stmt.Close()

	var updatedTodo Todo
	err = stmt.QueryRow(
		t.Title,
		t.Description,
		t.Status,
		t.CreatedBy,
		t.ID,
	).Scan(
		&updatedTodo.ID,
		&updatedTodo.Title,
		&updatedTodo.Description,
		&updatedTodo.Status,
	)
	if err != nil {
		return Todo{}, err
	}

	return updatedTodo, nil
}

func (t *Todo) DeleteTodo() error {

	query := `DELETE FROM todos
		WHERE created_by = ? AND id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(t.CreatedBy, t.ID)
	if err != nil {
		return err
	}

	if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}

func ConvertDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return dt.In(loc).Format(time.RFC822Z)
}
