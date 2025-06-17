package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Users	UserModel
	Tokens 	TokenModel
	Todos 	TodoModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Todos: TodoModel{DB: db},
	}
}