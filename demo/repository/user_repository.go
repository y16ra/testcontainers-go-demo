package repository

import (
	"database/sql"

	"github.com/y16ra/testcontainers-go-demo/demo/model"
)

type UserRepository interface {
	FindById(id int64) (*model.User, error)
	Store(user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindById(id int64) (*model.User, error) {
	user := &model.User{}
	if err := r.db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Store(user *model.User) error {
	err := r.db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", user.Name).Scan(&user.ID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
