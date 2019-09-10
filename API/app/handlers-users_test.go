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

func TestSearchUsers(t *testing.T) {

	//set up router and server
	router := hr.New()
	s := Server{DB: &mockDB{}, Router: router}
	s.Routes()

	//set up recorder and request
	rr := httptest.NewRecorder()

	url := fmt.Sprintf("/api/search?q=User-1")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	//run
	router.ServeHTTP(rr, req)

	//test status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code:\ngot: %v\n want: %v", status, http.StatusOK)
	}

	//test content type
	expContentType := "application/json"
	if cType := rr.Header().Get("Content-Type"); cType != expContentType {
		t.Errorf("content type header does not match:\ngot: %v\nwant: %v",
			cType, expContentType)
	}

	//set expectations for what the handler should return
	want := []*models.User{}
	User1 := &models.User{ID: userID, Name: "User-1", Avatar: "sailboat.jpg"}
	want = append(want, User1)

	//check what is actually returned
	got := []*models.User{}
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	//compare results
	same, gw := compareUserList(got, want)

	if !same {
		t.Error(gw)
	}

}

//compareUserList compares got vs want for a collection of posts
func compareUserList(got, want []*models.User) (bool, string) {
	for key, PostGot := range got {
		if PostGot.ID != want[key].ID {
			gw := fmt.Sprintf("handler returned wrong User ID:\ngot: %v\nwant: %v", PostGot.ID, want[key].ID)
			return false, gw
		}
		if PostGot.Name != want[key].Name {
			gw := fmt.Sprintf("handler returned wrong User Name:\ngot: %v\nwant: %v", PostGot.Name, want[key].Name)
			return false, gw
		}
		if PostGot.Avatar != want[key].Avatar {
			gw := fmt.Sprintf("handler returned wrong User Avatar:\ngot: %v\nwant: %v", PostGot.Avatar, want[key].Avatar)
			return false, gw
		}
	}
	return true, ""
}
