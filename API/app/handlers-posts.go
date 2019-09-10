package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/chiips/snippets/API/models"
	hr "github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

//allPosts retrieves all users' posts from the database.
func (s *Server) allPosts() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ hr.Params) {

		ctx := r.Context()

		//set limit of 10 results
		limit := 10

		//default prevDate is now to start at most recent result.
		prevDate := time.Now().UTC().Format(time.RFC3339)

		//if previous date is in the query then update prevDate to get results from that point in time back.
		prevDateQuery, ok := r.URL.Query()["prev"]

		if ok && len(prevDateQuery[0]) >= 1 && strings.TrimSpace(prevDateQuery[0]) != "" {
			prevDate = prevDateQuery[0]
		}

		////create a postsCh to communicate results and an error channe to communicate errors
		postsCh := make(chan []*models.Post)
		errCh := make(chan error)

		//send a separate goroutine to search the database.
		go func() {

			//check if the request context is cancelled by the time we get to here.
			if ctx.Err() != nil {
				return
			}

			//call the database
			posts, err := s.DB.AllPosts(prevDate, limit)

			//check if the request context is cancelled by the time we're done searching the database.
			if ctx.Err() != nil {
				return
			}

			//if the database search returns an error then send an error back on the error channel.
			if err != nil {
				errCh <- err
				return
			}

			//if the database search returns post results successfully then send the results back on the post channel.
			postsCh <- posts
			return

		}()

		//listen for three options in our program:
		select {
		//1. the context was cancelled.
		case <-ctx.Done():
			s.Log.Errorln(ctx.Err())
			http.Error(w, "We could not process your request at this time. Please try again later.", http.StatusRequestTimeout)
			return
		//2. there was an error searching the database.
		case err := <-errCh:
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		//3. success
		case posts := <-postsCh:
			//send the results as JSON data to the client
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(posts)
			if err != nil {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			return
		}
	}
}

//submitPost handles users submitting a new post
func (s *Server) submitPost() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ hr.Params) {

		ctx := r.Context()

		//confirm that the user id is present.
		//the user id should be passed into the context in the authenticateJWT middleware.
		currentUser, ok := ctx.Value(userContextKey).(uuid.UUID)
		if !ok {
			s.Log.Errorln("no userID in context")
			http.Error(w, http.StatusText(500), http.StatusForbidden)
			return
		}

		//confirm that the user id is not nil
		if uuid.Equal(currentUser, uuid.Nil) {
			s.Log.Errorln("userID came in with nil value.")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//get the information of the new submission
		submission := models.Post{}
		err = json.NewDecoder(r.Body).Decode(&submission)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		title := submission.Title
		body := submission.Body

		//check that the necessary submission information is present and properly formatted
		if strings.TrimSpace(title) == "" || strings.TrimSpace(body) == "" {
			s.Log.Errorln("invalid title and/or body")
			http.Error(w, "invalid title and/or body", http.StatusBadRequest)
			return
		}

		//check that the title and body are not too many characters (rather than checking for too many bytes)
		if utf8.RuneCountInString(title) > 50 || utf8.RuneCountInString(body) > 5000 {
			s.Log.Errorln("invalid title and/or body")
			http.Error(w, "invalid title and/or body", http.StatusBadRequest)
			return
		}

		//create uuid for post
		id, err := uuid.NewV4()
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		//build post
		post := &models.Post{}
		post.ID = id
		post.Title = title
		post.Body = body
		post.Created = time.Now().UTC()
		post.Updated = time.Now().UTC()
		post.Author.ID = currentUser

		//create success channel and error channel.
		okCh := make(chan bool)
		errCh := make(chan error)

		go func() {

			if ctx.Err() != nil {
				return
			}

			err = s.DB.CreatePost(post)

			//if the post request successfully reached the database then the submission was successful.

			if err != nil {
				errCh <- err
				return
			}

			okCh <- true
			return

		}()

		select {
		case <-ctx.Done():
			s.Log.Errorln(ctx.Err())
			http.Error(w, "We could not process your request at this time. Please try again later.", http.StatusRequestTimeout)
			return
		case err := <-errCh:
			s.Log.Errorln("error submitting post:", err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		case <-okCh:
			fmt.Fprint(w, "post submitted!")
			return
		}
	}

}

//editPost handles when users modify their submissions
func (s *Server) editPost() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps hr.Params) {

		ctx := r.Context()

		currentUser, ok := ctx.Value(userContextKey).(uuid.UUID)
		if !ok {
			s.Log.Errorln("no userID in context")
			http.Error(w, http.StatusText(500), http.StatusForbidden)
			return
		}

		if uuid.Equal(currentUser, uuid.Nil) {
			s.Log.Errorln("userID came in with nil value.")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		submission := &models.Post{}
		err = json.NewDecoder(r.Body).Decode(&submission)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		if strings.TrimSpace(submission.Title) == "" || strings.TrimSpace(submission.Body) == "" {
			s.Log.Errorln("invalid title and/or body")
			http.Error(w, "invalid title and/or body", http.StatusBadRequest)
			return
		}

		if utf8.RuneCountInString(submission.Title) > 50 || utf8.RuneCountInString(submission.Body) > 5000 {
			s.Log.Errorln("invalid title and/or body")
			http.Error(w, "invalid title and/or body", http.StatusBadRequest)
			return
		}

		//confirm the current user is the author of the submission
		if submission.Author.ID != currentUser {
			s.Log.Errorln("forbidden request")
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}

		//change last updated to now
		submission.Updated = time.Now().UTC()

		okCh := make(chan bool)
		errCh := make(chan error)

		go func() {

			if ctx.Err() != nil {
				return
			}

			err = s.DB.UpdatePost(submission)

			if err != nil {
				errCh <- err
				return
			}

			okCh <- true
			return

		}()

		select {
		case <-ctx.Done():
			s.Log.Errorln(ctx.Err())
			http.Error(w, "We could not process your request at this time. Please try again later.", http.StatusRequestTimeout)
			return
		case err := <-errCh:
			s.Log.Errorln("error editing post:", err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		case <-okCh:
			fmt.Fprint(w, "post edited!")
			return
		}

	}

}

//deletePost handles requests to remove one specific post
func (s *Server) deletePost() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps hr.Params) {

		ctx := r.Context()

		currentUser, ok := ctx.Value(userContextKey).(uuid.UUID)
		if !ok {
			s.Log.Errorln("no userID in context")
			http.Error(w, http.StatusText(500), http.StatusForbidden)
			return
		}

		if uuid.Equal(currentUser, uuid.Nil) {
			s.Log.Errorln("userID came in with nil value.")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		urlID := ps.ByName("postid")

		if urlID == "" {
			s.Log.Errorln("postid came in with zero value")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		id, err := uuid.FromString(urlID)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//query the database for the full post information.
		post, err := s.DB.OnePost(id)
		switch {
		case err == sql.ErrNoRows:
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		case err != nil:
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		//confirm post belongs to the current user
		if post.Author.ID != currentUser {
			s.Log.Errorln("forbidden request")
			http.Error(w, http.StatusText(500), http.StatusForbidden)
			return
		}

		okCh := make(chan bool)
		errCh := make(chan error)

		go func() {

			if ctx.Err() != nil {
				return
			}

			err = s.DB.DeletePost(post)

			if err != nil {
				errCh <- err
				return
			}

			okCh <- true
			return

		}()

		select {
		case <-ctx.Done():
			s.Log.Errorln(ctx.Err())
			http.Error(w, "We could not process your request at this time. Please try again later.", http.StatusRequestTimeout)
			return
		case err := <-errCh:
			s.Log.Errorln("error deleting post:", err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		case <-okCh:
			fmt.Fprint(w, "post deleted!")
			return
		}
	}
}
