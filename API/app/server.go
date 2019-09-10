package app

import (
	"github.com/chiips/snippets/API/logs"
	"github.com/chiips/snippets/API/models"
	hr "github.com/julienschmidt/httprouter"
)

//Server struct includes our datastore, router, and logger.
//All handlers hang off this Server struct to access its components via dependency injection as needed.
type Server struct {
	DB     models.Datastore
	Router *hr.Router
	Log    *logs.Log
}
