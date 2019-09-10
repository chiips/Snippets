package app

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/chiips/snippets/API/models"
	hr "github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//searchUsers checks the query parameter of the request and returns 10 results at a time
func (s *Server) searchUsers() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ hr.Params) {
		ctx := r.Context()

		//set limit of 10 results
		limit := 10

		//default prevDate is now to start at most recent result.
		prevDate := time.Now().UTC().Format(time.RFC3339)

		//if previous date is in the query then update prevDate to get results from that point in time back.
		prevDateQuery, ok := r.URL.Query()["prev"]

		if ok && strings.TrimSpace(prevDateQuery[0]) != "" && len(prevDateQuery[0]) >= 1 {
			prevDate = prevDateQuery[0]
		}

		//check the query
		var query string

		searchQuery, ok := r.URL.Query()["q"]

		if ok && strings.TrimSpace(searchQuery[0]) != "" && len(searchQuery[0]) >= 1 {
			query = searchQuery[0]
		} else {
			s.Log.Errorln("invalid search query")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//escape Postgres wildcard characters
		if strings.Contains(query, "%") {
			query = strings.Replace(query, "%", "\\%", -1)
		}

		if strings.Contains(query, "_") {
			query = strings.Replace(query, "_", "\\_", -1)
		}

		//create a usersCh to communicate results and an error channe to communicate errors
		usersCh := make(chan []*models.User)
		errCh := make(chan error)

		//send a separate goroutine to search the database.
		go func() {

			//check if the request context is cancelled by the time we get to here.
			if ctx.Err() != nil {
				return
			}

			//call the database
			users, err := s.DB.SearchUsers(query, prevDate, limit)

			//check if the request context is cancelled by the time we're done searching the database.
			if ctx.Err() != nil {
				return
			}

			//if the database search returns an error then send an error back on the error channel.
			if err != nil {
				errCh <- err
				return
			}

			//if the database search returns user results successfully then send the results back on the user channel.
			usersCh <- users
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
		case users := <-usersCh:
			//send the results as JSON data to the client
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(users)
			if err != nil {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			return
		}
	}
}

//signup checks a new account request and logs the user in
//login handler would perform similar checks and send tokens in cookies. logout handler would delete those cookies.
func (s *Server) signup() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ hr.Params) {

		ctx := r.Context()

		//get the information of the new account request
		signingUp := models.User{}
		err = json.NewDecoder(r.Body).Decode(&signingUp)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		name := signingUp.Name
		email := signingUp.Email
		password := signingUp.Password

		//check that the necessary account information is present and properly formatted
		if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
			s.Log.Errorln("bad form request")
			http.Error(w, "invalid name, email, and/or password", http.StatusBadRequest)
			return
		}

		//check name format
		rxName := regexp.MustCompile("^[a-zA-Z0-9_]{1,15}$")

		if !rxName.MatchString(name) {
			s.Log.Errorln("invalid name")
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}

		//check email format
		rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		//email addresses cannot be more than 254 bytes
		if len(email) > 254 || !rxEmail.MatchString(email) {
			s.Log.Errorln("invalid email address")
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}

		//check password format and length using helper function
		if !passwordIsValid(password) {
			s.Log.Errorln("invalid password")
			http.Error(w, "invalid password", http.StatusBadRequest)
			return
		}

		//check that the email is not already in use
		exists, err := s.DB.EmailCheck(email)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		if exists {
			s.Log.Errorln("email already taken")
			http.Error(w, "there is already an account with that email address", http.StatusBadRequest)
			return
		}

		//check that the name is not already taken
		exists, err = s.DB.NameCheck(name)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		if exists {
			s.Log.Errorln("username already taken")
			http.Error(w, "username already taken", http.StatusBadRequest)
			return
		}

		//hash the password
		bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		pwd := string(bs)

		//create uuid for user
		id, err := uuid.NewV4()
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		//build user
		user := &models.User{}
		user.ID = id
		user.Name = name
		user.Email = email
		user.Password = pwd
		user.Avatar = "puppy.jpg" //default photo for all new users
		user.Created = time.Now().UTC()
		user.Updated = time.Now().UTC()

		//create a success channel and an error channel
		okCh := make(chan bool)
		errCh := make(chan error)

		go func() {

			if ctx.Err() != nil {
				return
			}

			err = s.DB.CreateUser(user)

			//no checking context again here.
			//if the post request successfully reached the database then the user has an account.
			//the use won't get logged in but will be able to.

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
			s.Log.Errorln("error signing up:", err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		case <-okCh:
			//create a JWT, split into headerpaylod and signature, and put each into cookies.
			headerpayload, signature, err := s.createJWT(user.ID)
			if err != nil {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			//header and payload in non-HttpOnly cookie
			c1 := &http.Cookie{
				Name:     "token-hp",
				Value:    headerpayload,
				Secure:   true, //for testing over http, set Secure to false.
				Path:     "/",
				MaxAge:   0,
				SameSite: http.SameSiteDefaultMode,
			}
			http.SetCookie(w, c1)
			//signature in HttpOnly cookie
			c2 := &http.Cookie{
				Name:     "token-s",
				Value:    signature,
				Secure:   true,
				HttpOnly: true,
				Path:     "/",
				MaxAge:   0,
				SameSite: http.SameSiteDefaultMode,
			}
			http.SetCookie(w, c2)
			fmt.Fprint(w, "account created!")
			return
		}

	}
}

//passwordIsValid checks that the proposed password meets the following criteria.
func passwordIsValid(s string) bool {

	var minLen, upper, lower, number, special bool

	if len(s) >= 8 {
		minLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsLower(char):
			lower = true
		case unicode.IsNumber(char):
			number = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char): //spaces are allowed
			special = true
		}
	}
	return minLen && upper && lower && number && special
}

//set the maximum upload size of the image: 1MB
const maxUploadSize = 1024 * 1024

//editProfilePhoto handles a user uploading a new avatar image.
func (s *Server) editProfilePhoto() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps hr.Params) {

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

		//get id from url
		urlID := ps.ByName("userid")

		//confirm url id is not empty.
		if urlID == "" {
			s.Log.Errorln("userid came in with zero value.")
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}

		//convert id from url to type uuid.UUID for comparison.
		id, err := uuid.FromString(urlID)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		//confirm currentUser equals url id. if not then this request is forbidden.
		if !uuid.Equal(currentUser, id) {
			s.Log.Errorln("forbidden request.")
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}

		//start building user
		user := &models.User{}
		user.ID = id

		//validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			s.Log.Errorln("avatar file upload too big (>1MB)")
			http.Error(w, "file too big (>1MB)", http.StatusBadRequest)
			return
		}

		//get the image file
		mf, _, err := r.FormFile("avatar")

		if err != nil {
			s.Log.Errorln(err)
			if err == http.ErrMissingFile {
				http.Error(w, "missing file", http.StatusBadRequest)
			} else {
				http.Error(w, "invalid file", http.StatusBadRequest)
			}
			return
		}
		defer mf.Close()

		//make buffer to read the file header (512 bytes)
		buff := make([]byte, 512)
		if _, err = mf.Read(buff); err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		//check that the content type of the file is jpg or png
		fileType := http.DetectContentType(buff)

		switch fileType {
		case "image/jpeg", "image/jpg":
		case "image/x-png", "image/png":
			break
		default:
			s.Log.Errorln("bad form request: avatar of invalid file type")
			http.Error(w, "avatar of invalid file type", http.StatusBadRequest)
			return
		}

		//make random file name
		b := make([]byte, 12)
		rand.Read(b)
		fileName := fmt.Sprintf("%x", b)

		//get extension based on file type (includes the ".")
		fileEndings, err := mime.ExtensionsByType(fileType)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, "cannot read file type", http.StatusInternalServerError)
			return
		}

		//if got to here then the request was good

		//build full file name
		avatarName := fileName + fileEndings[0]

		//build user
		user.Avatar = avatarName
		user.Updated = time.Now().UTC()

		//go add to file system and database

		//create two done channels, one for each task.
		doneCh := make(chan bool, 2)
		errCh := make(chan error)

		//set handler context with cancel
		handlerCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {

			//check cancelled request.
			if ctx.Err() != nil {
				return
			}

			//check if the second task fails.
			if handlerCtx.Err() != nil {
				return
			}

			//create and store the file on the server
			err := createFile(mf, avatarName, currentUser)

			//if there's an error then also cancel the context to halt the second task.
			if err != nil {
				//send the error on the error channel first to send an error to the client instead of cancelled message.
				errCh <- err
				cancel()
				return
			}

			doneCh <- true
			return

		}()

		go func() {

			//check cancelled request
			if ctx.Err() != nil {
				return
			}

			//check if the second task fails.
			if handlerCtx.Err() != nil {
				return
			}

			err = s.DB.UpdateUserPhoto(user)

			//if there's an error then also cancel the context to halt the first task.
			if err != nil {
				errCh <- err
				cancel()
				return
			}

			doneCh <- true
			return

		}()

		//listen for three options:
		doneTasks := 0
		for {
			select {
			//1. context cancelled
			case <-ctx.Done():
				s.Log.Errorln(ctx.Err())
				http.Error(w, "We could not process your request at this time. Please try again later.", http.StatusRequestTimeout)
				return
			//2. error occurred
			case err := <-errCh:
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			//3. both tasks were completed successfully
			case <-doneCh:
				doneTasks++
				if doneTasks == 2 {
					w.Header().Set("Content-Type", "application/json")
					err = json.NewEncoder(w).Encode(user)
					if err != nil {
						s.Log.Errorln(err)
						http.Error(w, http.StatusText(500), http.StatusInternalServerError)
						return
					}
					return
				}
			}
		}

	}

}

//createFile stores the incoming image on the server.
func createFile(mf multipart.File, avatarName string, currentUser uuid.UUID) error {
	//get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	//set a folder path including the id of the current user. each user gets their own folder.
	folderPath := fmt.Sprintf("private/assets/%s", currentUser)
	//check the status of the folder path. if it doesn't exist then create it.
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		//the user already has a path, i.e., they've already uploaded a unique avatar before.
		oldPath := filepath.Join(wd, folderPath)

		//remove the contents of existing path using a helper function
		err := removeContents(oldPath)
		if err != nil {
			return err
		}

	}

	//create the new file path
	newPath := filepath.Join(wd, folderPath, avatarName)

	nf, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer nf.Close()

	//copy the image into our file
	mf.Seek(0, 0)
	io.Copy(nf, mf)
	return nil
}

//removeContents clears all files in the given directory.
func removeContents(dir string) error {
	//open the directory
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	//read all the files in the directory
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	//range over the files and remove their paths
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//deleteUser handles account deletions
func (s *Server) deleteUser() hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps hr.Params) {

		ctx := r.Context()

		currentUser, ok := ctx.Value(userContextKey).(uuid.UUID)
		if !ok {
			s.Log.Errorln("no userID in context")
			http.Error(w, http.StatusText(500), http.StatusForbidden)
			return
		}

		urlID := ps.ByName("userid")

		if urlID == "" {
			s.Log.Errorln("userid came in with zero value.")
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}

		id, err := uuid.FromString(urlID)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		if !uuid.Equal(currentUser, id) {
			s.Log.Errorln("forbidden request.")
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}

		//delete user's photos
		wd, err := os.Getwd()
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		folderPath := fmt.Sprintf("private/assets/%s", currentUser)

		//set the path. if it exists, i.e., they've uploaded a unique avatar, then find and delete their folder.
		if _, err := os.Stat(folderPath); !os.IsNotExist(err) {

			userPath := filepath.Join(wd, folderPath)

			//first remove contents
			err := removeContents(userPath)
			if err != nil {
				s.Log.Errorln(err)
			}

			//then remove path
			err = os.Remove(userPath)
			if err != nil {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

		}

		//delete the user from the database
		user := &models.User{}
		user.ID = id

		err = s.DB.DeleteUser(user)
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		//check the user's cookies
		//both cookies had to come in to verify the request, but always check errors.
		c1, err := r.Cookie("token-hp")
		if err != nil {
			//if there's an error then an unauthorized response will prompt the client to redirect the user anyway.
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		c2, err := r.Cookie("token-s")
		if err != nil {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		//remove the cookies from the browser
		c1 = &http.Cookie{
			Name:   "token-hp",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, c1)

		c2 = &http.Cookie{
			Name:   "token-s",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, c2)

		fmt.Fprint(w, "account successfully deleted.")
		return
	}
}
