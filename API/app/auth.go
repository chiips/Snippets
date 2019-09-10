package app

import (
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

//MyClaims struct defined for adding to jwt.StandardClaims as an embedded type.
type MyClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

//createJWT creates a new JWT for a user.
func (s *Server) createJWT(id uuid.UUID) (string, string, error) {

	//5 minute expiration time in unix milliseconds
	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	//Create the JWT claims which include the user id, expiry time, and issuer
	claims := &MyClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    os.Getenv("jwt_issuer"),
		},
	}

	//Declare the token with the algorithm used for signing and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string using the JWT key
	tokenString, err := token.SignedString([]byte(os.Getenv("jwt_key")))
	if err != nil {
		return "", "", err
	}

	//split token into header+paylod and signature for sending in two separate cookies.
	//headerpayload will go into a non-HttpOnly cookie for the front-end client to access via Javascript.
	//signature will go into an HttpOnly cookie for security.
	tokenSplit := strings.Split(tokenString, ".")
	header, payload, signature := tokenSplit[0], tokenSplit[1], tokenSplit[2]

	headerpaylod := header + "." + payload

	return string(headerpaylod), string(signature), nil
}
