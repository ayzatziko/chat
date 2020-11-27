package mem

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ayzatziko/chat/repository"
)

var _ repository.Users = (*Users)(nil)

type Users struct {
	mu         sync.Mutex
	byUsername map[string]*repository.User
	byNickname map[string]string
}

func New() *Users {
	return &Users{
		byUsername: make(map[string]*repository.User),
		byNickname: make(map[string]string),
	}
}

func (r *Users) UserByUsername(username string) (*repository.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.byUsername[username]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (r *Users) UserByNickname(nickname string) (*repository.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	un, ok := r.byNickname[nickname]
	if !ok {
		return nil, sql.ErrNoRows
	}

	u, ok := r.byUsername[un]
	if !ok {
		delete(r.byNickname, nickname)
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (r *Users) UpdateUser(u *repository.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	usr, ok := r.byUsername[u.Username]
	if !ok {
		return sql.ErrNoRows
	}

	usr.Nickname = u.Nickname
	return nil
}

func (r *Users) NicknameExists(nickname string) (bool, error) {
	_, ok := r.byNickname[nickname]
	return ok, nil
}

func (r *Users) InsertUser(u *repository.User) error {
	if _, ok := r.byUsername[u.Username]; ok {
		return fmt.Errorf("user with username %q already exists", u.Nickname)
	}
	if _, ok := r.byNickname[u.Nickname]; ok {
		return fmt.Errorf("user with nickname %q already exists", u.Nickname)
	}

	r.byUsername[u.Username] = u
	r.byNickname[u.Nickname] = u.Username
	return nil
}
