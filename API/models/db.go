package models

import (
	"database/sql"

	//pq is necessary for connecting with PostgreSQL
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

//Datastore is an interface to work with the Postgres database.
//The Server struct in API/app/server.go includes this Datastore interface for handlers to access via dependency injection.
//Using an interface allows us to easily create mock databases for testing purposes.
type Datastore interface {

	//Sample User methods
	SearchUsers(query, prevDate string, limit int) ([]*User, error)
	CreateUser(user *User) error
	EmailCheck(email string) (bool, error)
	NameCheck(name string) (bool, error)
	UpdateUserPhoto(user *User) error
	DeleteUser(user *User) error

	//Sample Post methods
	AllPosts(prevDate string, limit int) ([]*Post, error)
	OnePost(id uuid.UUID) (*Post, error)
	CreatePost(Post *Post) error
	UpdatePost(Post *Post) error
	DeletePost(Post *Post) error
}

//DB is our database type
//By attaching the Datastore interface's methods, our DB struct will implement the Datastore interface.
type DB struct {
	*sql.DB
}

//NewDB creates a new DB instance
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
