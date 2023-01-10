package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tieubaoca/go-chat-server/types"
)

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
		log.Error(err)
		return "", "", err
	}
	defer resp.Body.Close()
	var body map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	accessToken, ok := body["access_token"]
	if !ok {
		log.Error(types.ErrorInvalidInput)
		return "", "", errors.New(types.ErrorInvalidInput)
	}
	refreshToken, ok := body["refresh_token"]
	if !ok {
		log.Error(types.ErrorInvalidInput)
		return "", "", errors.New(types.ErrorInvalidInput)
	}
	return accessToken.(string), refreshToken.(string), nil

}
