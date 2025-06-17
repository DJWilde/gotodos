package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Todo struct {
	ID			int64		`json:"id"`
	CreatedAt	time.Time	`json:"-"`
	Name 		string 		`json:"name"`
	Description string		`json:"description"`
	Done 		bool 		`json:"done"`
	UserID 		int64		`json:"user_id"`
	Version		int32		`json:"version"`
}

type TodoModel struct {
	DB *sql.DB
}

func (m TodoModel) GetByID(id int64) (*Todo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, description, done, user_id, version
		FROM todos
		WHERE id = $1`

	var todo Todo

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Name,
		&todo.Description,
		&todo.Done,
		&todo.UserID,
		&todo.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &todo, nil
}

func (m TodoModel) GetByUserID(userID int64) ([]*Todo, error) {
	query := `
		SELECT id, created_at, name, description, done, user_id, version
		FROM todos
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []*Todo{}

	for rows.Next() {
		var todo Todo

		err := rows.Scan(
			&todo.ID,
			&todo.CreatedAt,
			&todo.Name,
			&todo.Description,
			&todo.Done,
			&todo.UserID,
			&todo.Version,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (m TodoModel) Insert(todo *Todo, userID int64) error {
	query := `
		INSERT INTO todos (name, description, done, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{todo.Name, todo.Description, todo.Done, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

func (m TodoModel) Update(todo *Todo) error {
	query := `
		UPDATE todos
		SET name = $1, description = $2, done = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []any{
		todo.Name,
		todo.Description,
		todo.Done,
		todo.ID,
		todo.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m TodoModel) DeleteByID(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM todos
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}