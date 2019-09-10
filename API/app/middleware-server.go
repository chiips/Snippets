package app

import (
	"context"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
)

//These middlewares protect the server's router and therefore apply to all routes.

//NewLimiter sets up tollbooth rate limiter
func (s *Server) NewLimiter() *limiter.Limiter {

	lmt := tollbooth.NewLimiter(2, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		log := s.Log.WithFields(log.Fields{"request id": r.Header.Get("X-REQUEST-ID"), "request uri": r.RequestURI, "request method": r.Method})
		log.Errorln("request limit reached")
	})

	return lmt

}

//Timeout sets a context withcancel that matches &http.Server read + write timeout in main.go.
//This middleware allows the handlers to respond precisely to timeout errors while the &http.Server timeout serves as an absolute safeguard.
func (s *Server) Timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		r = r.WithContext(tctx)

		next.ServeHTTP(w, r)

	})

}

//LogRequests logs the custom request ID (created in the SPA), URI, and method to the logger.
func (s *Server) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := s.Log.WithFields(log.Fields{"id": r.Header.Get("X-REQUEST-ID"), "uri": r.RequestURI, "method": r.Method})

		log.Infoln("about to serve")

		next.ServeHTTP(w, r)

		log.Infoln("finished serving")

	})

}

//SetHeaders sets the CSRF Token for all responses
func (s *Server) SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-CSRF-Token", csrf.Token(r))

		next.ServeHTTP(w, r)

	})

}

//CSRFErrorHandler is a custom error handler when CSRF tokens come in invalid
func (s *Server) CSRFErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s.Log.Errorln(csrf.FailureReason(r))
		http.Error(w, "CSRF token invalid", http.StatusForbidden)
		return

	})

}
