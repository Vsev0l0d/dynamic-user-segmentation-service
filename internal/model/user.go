package model

import (
	"fmt"
	"net/http"
)

type (
	User struct {
		Id int `json:"id" db:"id"`
	}

	UserList struct {
		Users []User `json:"users"`
	}
)

func (u *User) Bind(r *http.Request) error {
	if u.Id <= 0 {
		return fmt.Errorf("id is a required field that accepts values greater than zero")
	}
	return nil
}
func (*UserList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
