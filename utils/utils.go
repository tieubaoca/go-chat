package utils

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"
)

func GetAccessTokenByReq(r *http.Request) string {
	var tokenString string
	tokenCookie, err := r.Cookie("access-token")
	if err != nil {
		tokenString = r.URL.Query().Get("access-token")
	} else {
		tokenString = tokenCookie.Value
	}
	return tokenString
}

func ParseUnverified(tokenString string) (*jwt.Token, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return token, nil
}

func Parse(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		keyStr := os.Getenv("PUBKEY_HEADER") + "\n" + os.Getenv("JWT_PUBKEY") + "\n" + os.Getenv("PUBKEY_FOOTER")
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	testToken := jwt.New(jwt.SigningMethodEdDSA)
	testToken.Claims = token.Claims.(jwt.MapClaims)
	return token, nil
}

func GetSessionIdFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["sid"].(string)
}

func GetUsernameFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["preferred_username"].(string)
}

func GetSaIdFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["sub"].(string)
}
