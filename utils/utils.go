package utils

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/golang-jwt/jwt"
)

func GetAccessTokenByReq(r *http.Request) string {
	accessToken := getTokenFromHeader(r)
	if accessToken == "" {
		accessToken = getTokenFromQueryParams(r)
	}
	if accessToken == "" {
		accessToken = getTokenFromCookie(r)
	}
	return accessToken
}

func getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("access-token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getTokenFromQueryParams(r *http.Request) string {
	return r.URL.Query().Get("access-token")
}

func getTokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func JWTParseUnverified(tokenString string) (*jwt.Token, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return token, nil
}

func JWTSaasParse(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		keyStr := os.Getenv("PUBKEY_HEADER") + "\n" + os.Getenv("JWT_PUBKEY") + "\n" + os.Getenv("PUBKEY_FOOTER")
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
	})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	testToken := jwt.New(jwt.SigningMethodEdDSA)
	testToken.Claims = token.Claims.(jwt.MapClaims)
	return token, nil
}

func JWTGenerateToken(data interface{}) (string, error) {
	type jwtCustomClaim struct {
		Data interface{} `json:"data"`
		jwt.StandardClaims
	}
	claims := &jwtCustomClaim{
		data,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// var sampleSecretKey = []byte("SecretYouShouldHide")
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_PRIV_KEY")))
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", err
	}
	return tokenString, nil
}

func JWTVerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_PRIV_KEY")), nil
	})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return token, nil
}

func GetSessionIdFromToken(token *jwt.Token) string {
	return token.Claims.(jwt.MapClaims)["sid"].(string)
}

func GetMessageIds(messages []models.Message) []string {
	var ids []string
	for _, message := range messages {
		ids = append(ids, message.Id.Hex())
	}
	return ids
}
