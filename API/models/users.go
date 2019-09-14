package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//User type defined
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	Avatar   string    `json:"avatar"`
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
}

//Our selection of sample User methods to satisfy the Datastore interface:

//SearchUsers takes a search query and limit and returns all posts in reverse chronological order or an error.
func (db *DB) SearchUsers(query, prevDate string, limit int) ([]*User, error) {
	users := []*User{}

	rows, err := db.Query("SELECT id, name, avatar, created FROM users WHERE name ILIKE '%' || $1 || '%' AND created < $2 ORDER BY created DESC LIMIT $3", query, prevDate, limit)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Avatar, &user.Created)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

//CreateUser creates a new user and returns nil or an error
//CreateUser expects user will come in with name string, email string, pwd []byte
func (db *DB) CreateUser(user *User) error {

	_, err := db.Exec("INSERT INTO users (id, name, email, password, avatar, created, updated) VALUES ($1, $2, $3, $4, $5, $6, $7)", user.ID, user.Name, user.Email, user.Password, user.Avatar, user.Created, user.Updated)
	if err != nil {
		return err
	}

	return nil
}

//EmailCheck checks if an email is already in use when a new user signs up.
func (db *DB) EmailCheck(email string) (bool, error) {

	var exists bool

	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);", email)
	err := row.Scan(&exists)
	if err != nil {
		return exists, err
	}

	return exists, err
}

//NameCheck checks if a name is already in use when a new user signs up
func (db *DB) NameCheck(name string) (bool, error) {

	var exists bool

	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = $1);", name)
	err := row.Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, err
}

//UpdateUserPhoto updates a user's profile photo and returns nil or an error.
//UpdateUserPhoto expects user will come in with avatar string, updated time.Time
func (db *DB) UpdateUserPhoto(user *User) error {

	_, err := db.Exec("UPDATE users SET avatar=$2, updated=$3 WHERE id=$1;", user.ID, user.Avatar, user.Updated)
	if err != nil {
		return err
	}

	return nil

}

//DeleteUser deletes one specific user from DB, along with associated posts, and returns nil or an error.
//DeleteUser expects user will come in with id uuid.UUID
func (db *DB) DeleteUser(user *User) error {

	_, err := db.Exec("DELETE FROM posts WHERE uid=$1;", user.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM users WHERE id=$1;", user.ID)
	if err != nil {
		return err
	}

	return nil
}
