package pq

import (
	"database/sql"

	"github.com/ayzatziko/chat/repository"
)

var _ repository.Users = (*Users)(nil)

type Users struct {
	db *sql.DB
}

func New(db *sql.DB) *Users {
	return &Users{db: db}
}

func (r *Users) UserByUsername(username string) (*repository.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE username=?", username)
	u := repository.User{}
	if err := row.Scan(&u.Username, &u.Nickname, &u.PasswordHash); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Users) UserByNickname(nickname string) (*repository.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE nickname=?", nickname)
	u := repository.User{}
	if err := row.Scan(&u.Username, &u.Nickname, &u.PasswordHash); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Users) UpdateUser(u *repository.User) error {
	_, err := r.db.Exec("UPDATE users SET nickname=? WHERE username=?", u.Nickname, u.Username)
	return err
}

func (r *Users) NicknameExists(nickname string) (bool, error) {
	row := r.db.QueryRow("SELECT COUNT(nickname) FROM users WHERE nickname=?", nickname)
	var c int
	if err := row.Scan(&c); err != nil {
		return false, err
	}
	return c == 0, nil
}

func (r *Users) InsertUser(u *repository.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, nickname, passwordhash) VALUES (?, ?, ?)",
		u.Username,
		u.Nickname,
		u.PasswordHash)
	return err
}
