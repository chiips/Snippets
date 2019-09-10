package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//Post type defined
type Post struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Author  User      `json:"author"`
}

//Our selection of sample Post methods to satisfy the Dataface interface:

//AllPosts takes a previous date and limit and returns all posts in reverse chronological order or an error.
func (db *DB) AllPosts(prevDate string, limit int) ([]*Post, error) {
	posts := []*Post{}

	rows, err := db.Query("SELECT posts.ID, posts.title, posts.body, posts.created, posts.updated, users.id, users.name, users.avatar FROM posts INNER JOIN users ON posts.uid = users.id WHERE posts.created < $1 ORDER BY posts.created DESC LIMIT $2;", prevDate, limit)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		Post := &Post{}
		err := rows.Scan(&Post.ID, &Post.Title, &Post.Body, &Post.Created, &Post.Updated, &Post.Author.ID, &Post.Author.Name, &Post.Author.Avatar)
		if err != nil {
			return posts, err
		}
		posts = append(posts, Post)
	}
	if err := rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil
}


//OnePost returns one specific post or an error
func (db *DB) OnePost(id uuid.UUID) (*Post, error) {

	post := &Post{}

	row := db.QueryRow("SELECT posts.ID, posts.title, posts.body, posts.created, posts.updated, users.id, users.name, users.avatar FROM posts INNER JOIN users ON posts.uid = users.id WHERE posts.id = $1", id)

	err := row.Scan(&post.ID, &post.Title, &post.Body, &post.Created, &post.Updated, &post.Author.ID, &post.Author.Name, &post.Author.Avatar)
	if err != nil {
		return post, err
	}

	return post, nil
}


//CreatePost creates a new post in the DB and returns an error.
//CreatePost expects Post will come in with id uuid.UUID, title string, body string, created time.Time, uid uuid.UUID
func (db *DB) CreatePost(Post *Post) error {

	_, err := db.Exec("INSERT INTO posts (id, title, body, created, updated, uid) VALUES ($1, $2, $3, $4, $5, $6)", Post.ID, Post.Title, Post.Body, Post.Created, Post.Updated, Post.Author.ID)
	if err != nil {
		return err
	}

	return nil
}

//UpdatePost updates a specific Post in DB and returns an error.
//UpdatePost expects Post will come in with id uuid.UUID, title string, body string, updated time.Time
func (db *DB) UpdatePost(Post *Post) error {

	_, err := db.Exec("UPDATE posts SET title=$2, body=$3, updated=$4 WHERE id=$1;", Post.ID, Post.Title, Post.Body, Post.Updated)
	if err != nil {
		return err
	}

	return nil
}

//DeletePost deletes one specific Post from DB and returns an error.
//DeletePost expects Post will come in with id uuid.UUID
func (db *DB) DeletePost(Post *Post) error {

	_, err := db.Exec("DELETE FROM posts WHERE id=$1;", Post.ID)
	if err != nil {
		return err
	}
	return nil
}
