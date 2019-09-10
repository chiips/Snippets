package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chiips/snippets/API/models"
	hr "github.com/julienschmidt/httprouter"
)

func TestAllPosts(t *testing.T) {

	//initialize router and server
	router := hr.New()
	s := Server{DB: &mockDB{}, Router: router}
	s.Routes()

	//set up recorder and request
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	//serve the request
	router.ServeHTTP(rr, req)

	//check the status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code:\ngot:\n%v\n want:\n%v", status, http.StatusOK)
	}

	//check the content type
	expContentType := "application/json"
	if ctype := rr.Header().Get("Content-Type"); ctype != expContentType {
		t.Errorf("content type header does not match:\ngot: \n%v\nwant: \n%v",
			ctype, expContentType)
	}

	//set expectations for what the handler should return
	want := []*models.Post{}
	Post1 := &models.Post{ID: postID1, Title: "Post 1", Body: "Body 1", Created: now, Updated: now, Author: models.User{ID: userID, Name: "User-1", Avatar: "sailboat.jpg"}}
	Post2 := &models.Post{ID: postID2, Title: "Post 2", Body: "Body 2", Created: now, Updated: now, Author: models.User{ID: userID, Name: "User-1", Avatar: "sailboat.jpg"}}
	want = append(want, Post1)
	want = append(want, Post2)

	//check what is actually returned
	got := []*models.Post{}
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	//compare results
	same, gw := comparePostList(got, want)

	if !same {
		t.Error(gw)
	}

}

//comparePostList compares got vs want for a collection of posts
func comparePostList(got, want []*models.Post) (bool, string) {
	for key, PostGot := range got {
		if PostGot.ID != want[key].ID {
			gw := fmt.Sprintf("handler returned wrong Post ID:\ngot: %v\nwant: %v", PostGot.ID, want[key].ID)
			return false, gw
		}
		if PostGot.Title != want[key].Title {
			gw := fmt.Sprintf("handler returned wrong Post Title:\ngot: %v\nwant: %v", PostGot.Title, want[key].Title)
			return false, gw
		}
		if PostGot.Body != want[key].Body {
			gw := fmt.Sprintf("handler returned wrong Post Body:\ngot: %v\nwant: %v", PostGot.Body, want[key].Body)
			return false, gw
		}
		if PostGot.Created != want[key].Created {
			gw := fmt.Sprintf("handler returned wrong Post Created:\ngot: %v\nwant: %v", PostGot.Created, want[key].Created)
			return false, gw
		}
		if PostGot.Updated != want[key].Updated {
			gw := fmt.Sprintf("handler returned wrong Post Updated:\ngot: %v\nwant: %v", PostGot.Updated, want[key].Updated)
			return false, gw
		}
		if PostGot.Author.ID != want[key].Author.ID {
			gw := fmt.Sprintf("handler returned wrong Post Author ID:\ngot: %v\nwant: %v", PostGot.Author.ID, want[key].Author.ID)
			return false, gw
		}
		if PostGot.Author.Name != want[key].Author.Name {
			gw := fmt.Sprintf("handler returned wrong Post Author Name:\ngot: %v\nwant: %v", PostGot.Author.Name, want[key].Author.Name)
			return false, gw
		}
		if PostGot.Author.Avatar != want[key].Author.Avatar {
			gw := fmt.Sprintf("handler returned wrong Post Author Avatar:\ngot: %v\nwant: %v", PostGot.Author.Avatar, want[key].Author.Avatar)
			return false, gw
		}
	}
	return true, ""
}
