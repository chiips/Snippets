package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/chiips/snippets/API/models"
	uuid "github.com/satori/go.uuid"
)

//generate variables for sample user and posts
var userID uuid.UUID
var postID1 uuid.UUID
var postID2 uuid.UUID
var now time.Time
var err error

func init() {

	userID, err = uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}

	postID1, err = uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}

	postID2, err = uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}

	now = time.Now().UTC()
}

//mockDB is a struct to mock our database for testing
type mockDB struct {
	//by embedding models.Datastore, mockDB implements the interface.
	//this way we don't need to stub each datastore method
	models.Datastore
}

//Sample user database method

func (mdb *mockDB) SearchUsers(query, prevDate string, limit int) ([]*models.User, error) {
	//create slice to serve as database.
	users := []*models.User{}
	users = append(users, &models.User{ID: userID, Name: "User-1", Avatar: "sailboat.jpg"})
	users = append(users, &models.User{ID: userID, Name: "User-2", Avatar: "sailboat.jpg"})
	users = append(users, &models.User{ID: userID, Name: "User-3", Avatar: "sailboat.jpg"})

	results := []*models.User{}

	for _, user := range users {
		if strings.Contains(user.Name, query) {
			results = append(results, user)
		}
	}

	return results, nil
}

//Sample post database method
func (mdb *mockDB) AllPosts(prevDate string, limit int) ([]*models.Post, error) {
	posts := []*models.Post{}
	posts = append(posts, &models.Post{ID: postID1, Title: "Post 1", Body: "Body 1", Created: now, Updated: now, Author: models.User{ID: userID, Name: "User-1", Avatar: "sailboat.jpg"}})
	posts = append(posts, &models.Post{ID: postID2, Title: "Post 2", Body: "Body 2", Created: now, Updated: now, Author: models.User{ID: userID, Name: "User-1", Avatar: "sailboat.jpg"}})
	return posts, nil
}
