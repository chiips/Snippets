package main

import (
	"fmt"
	"time"

	"net/http"
	"os"

	"github.com/chiips/snippets/API/app"
	"github.com/chiips/snippets/API/logs"
	"github.com/chiips/snippets/API/models"
	"github.com/didip/tollbooth"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	hr "github.com/julienschmidt/httprouter"
)

func main() {

	//load .env file for environment variables
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Println("env err:", envErr)
	}

	//set up logger
	logfile := os.Getenv("logfile")
	logger, err := logs.NewLogger(logfile)
	if err != nil {
		fmt.Println("error setting up new logger:", err)
	}

	//get database variables
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbHost := os.Getenv("db_host")
	dbName := os.Getenv("db_name")
	port := os.Getenv("server_port")

	//initiate Postgres connection with database variables
	dbURI := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, dbHost, dbName)
	db, err := models.NewDB(dbURI)
	if err != nil {
		logger.Panic(err)
	}
	//defer close db
	defer db.Close()

	//set up new router using Julien Schmidt's httprouter
	router := hr.New()

	//assign database, router, and logger to our app's Server struct
	s := app.Server{DB: db, Router: router, Log: logger}
	//initialize the Server's routes
	s.Routes()

	//Initiate CSRF protection
	key := []byte(os.Getenv("32-byte-auth-key"))
	errHandler := csrf.ErrorHandler(s.CSRFErrorHandler())
	security := csrf.Secure(false) //for development over http instead of https. true by default
	csrfProtect := csrf.Protect(key, errHandler, security)

	//initialize new limiter
	lmt := s.NewLimiter()

	//set our server object for ListenAndServe with all the server middleware
	srvHandler := s.Timeout(csrfProtect(s.LogRequests(s.SetHeaders(s.Router))))
	limitedHandler := tollbooth.LimitHandler(lmt, srvHandler)
	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
		Handler:      limitedHandler,
	}

	//listen and serve
	s.Log.Infoln("Listening on:", port)
	s.Log.Fatalln(srv.ListenAndServe()) //for development over HTTP instead of HTTPS
	//srv.ListenAndServeTLS()
}
