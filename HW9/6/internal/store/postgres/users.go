package postgres

import (
	"context"
	"fmt"
	"6/internal/models"
	"6/internal/store"
	"github.com/jmoiron/sqlx"

)

func (db *DB) Users() store.UsersRepository {
	if db.users == nil {
		db.users = NewUsersRepository(db.conn)
	}

	return db.users
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UsersRepository {
	return &UsersRepository{conn: conn}
}

func (u UsersRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.conn.Exec("INSERT INTO users(nick, password, bio, email) VALUES ($1, $2, %3, $4)", user.Nick, user.Password, user.Bio, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u UsersRepository) All(ctx context.Context, filter *models.UsersFilter) ([]*models.User, error) {
	users := make([]*models.User, 0)
	basicQuery := "SELECT * FROM users"
	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := u.conn.Select(&users, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return users, nil
	}

	if err := u.conn.Select(&users, basicQuery); err != nil {
		return nil, err
	}


	return users, nil
}

func (u UsersRepository) ByID(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.Get(user, "SELECT id, nick, password, bio, email FROM users WHERE id=$1", id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u UsersRepository) Update(ctx context.Context, user *models.User) error {
	_, err := u.conn.Exec("UPDATE users SET nick = $1, password = $2, bio = $3, email = $4 WHERE id = $5", user.Nick, user.Password, user.Bio, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u UsersRepository) Delete(ctx context.Context, id int) error {
	_, err := u.conn.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}