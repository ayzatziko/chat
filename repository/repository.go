package repository

import "fmt"

type User struct {
	Username     string
	Nickname     string
	PasswordHash string
}

type Users interface {
	UserByUsername(string) (*User, error)
	UserByNickname(string) (*User, error)
	UpdateUser(*User) error // should be just update nickname
	NicknameExists(string) (bool, error)
	InsertUser(*User) error
}

type Message struct {
	Author    string `json:"author"`
	Body      string `json:"body"`
	Recipient string `json:"recipient,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("<Message: <author: %s>, <body: %s>, <recipient: %s>>", m.Author, m.Body, m.Recipient)
}

type Messages interface {
	Add(*Message)
	Messages() []*Message
}
