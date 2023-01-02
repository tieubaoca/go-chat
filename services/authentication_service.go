package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func Parse(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		keyStr := os.Getenv("PUBKEY_HEADER") + "\n" + os.Getenv("JWT_PUBKEY") + "\n" + os.Getenv("PUBKEY_FOOTER")
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return token, nil
}

func GetAccessToken(username string, password string) (string, string, error) {
	req := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "keycloak.summonersarena.io",
			Path:   "/realms/summonersarena/protocol/openid-connect/token",
		},
		Header: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
		},
		Body: ioutil.NopCloser(
			strings.NewReader(
				"username=" + username +
					"&password=" + password +
					"&grant_type=password&client_id=" + os.Getenv("CLIENT_ID") +
					"&client_secret=" + os.Getenv("CLIENT_SECRET") + "")),
	}

	client := &http.Client{}

	resp, err := client.Do(&req)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	defer resp.Body.Close()
	var body map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	accessToken, ok := body["access_token"]
	refreshToken, _ := body["access_token"]
	if !ok {
		log.Println("No access token")
		return "", "", nil
	}
	return accessToken.(string), refreshToken.(string), nil

}
