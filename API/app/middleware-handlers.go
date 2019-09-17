package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	hr "github.com/julienschmidt/httprouter"
)

//define a custom contextkey type
type contextKey string

const userContextKey contextKey = "userID"

//authenticateJWT controls access to handlers according to the validity of the incoming JWT.
func (s *Server) authenticateJWT(next hr.Handle) hr.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps hr.Params) {

		//get the JWT header.payload
		c1, err := r.Cookie("token-hp")
		if err != nil {
			if err == http.ErrNoCookie {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//get the JWT signature
		c2, err := r.Cookie("token-s")
		if err != nil {
			if err == http.ErrNoCookie {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//combine for full JWT
		tknStr := c1.Value + "." + c2.Value

		//Parse the JWT string and store the result in `&MyClaims{}`.
		tkn, err := jwt.ParseWithClaims(tknStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("jwt_key")), nil
		})

		//check validity of the token
		if !tkn.Valid {
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}

		//catch any errors that might be missed even if the token is valid.
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				s.Log.Errorln(err)
				http.Error(w, http.StatusText(401), http.StatusUnauthorized)
				return
			}
			s.Log.Errorln(err)
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//make sure we can get the claims
		claims, ok := tkn.Claims.(*MyClaims)
		if !ok {
			s.Log.Errorln("invalid claims")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//verify the issuer field of the JWT
		verifiedIssuer := claims.VerifyIssuer(os.Getenv("jwt_issuer"), true)
		if !verifiedIssuer {
			s.Log.Errorln("invalid issuer")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		//reject if authenticated but trying to reach login
		requestPath := r.URL.Path
		requestPath = hr.CleanPath(requestPath)

		matched, err := regexp.MatchString("/api/login", requestPath)

		if matched && err == nil {
			s.Log.Errorln("login attempt with existing valid JWT")
			//status ok so frontend knows just redirect to logged in homepage
			w.WriteHeader(http.StatusOK)
			return
		}

		//put ID in context to pass along request chain
		ctx := context.WithValue(r.Context(), userContextKey, claims.ID)
		r = r.WithContext(ctx)

		s.Log.Infoln("JWT authentication OK, serving next")
		next(w, r, ps)

	}

}
