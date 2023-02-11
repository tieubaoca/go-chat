package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"github.com/tieubaoca/go-chat-server/utils/saasApi"

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

func GetCitizen(token string) (map[string]interface{}, error) {
	header := http.Header{}
	header.Add(
		"Authorization",
		"Bearer "+token,
	)
	res, err := saasApi.Get(
		os.Getenv("SAAS_HOST")+"/saas/api/v1/citizen/get-current-citizen",
		header,
		url.Values{},
	)
	if err != nil {
		log.ErrorLogger.Println(err)
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	if res.StatusCode != 200 {
		log.ErrorLogger.Println(result["message"])
		return nil, errors.New("get citizen error, status: " + res.Status)
	}
	data := result["data"].(map[string]interface{})
	log.InfoLogger.Println(data)
	return data, nil
}

func GetListFriendInfo(saasAccessToken string, paginationReq request.PaginationOnlineFriendReq) ([]interface{}, error) {

	header := http.Header{}
	header.Add(
		"Authorization",
		"Bearer "+saasAccessToken,
	)
	query := url.Values{
		"page": []string{fmt.Sprint(paginationReq.Page)},
		"size": []string{fmt.Sprint(paginationReq.Size)},
	}

	resp, err := saasApi.Get(
		os.Getenv("SAAS_HOST")+"/saas/api/v1/friend/getListFriendInfo",
		header,
		query,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	var resBody response.SaasResponse
	json.NewDecoder(resp.Body).Decode(&resBody)
	log.InfoLogger.Println(resBody)
	data := resBody.Data.([]interface{})
	return data, nil
}

func GetAllFriends(token string) ([]interface{}, error) {

	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)
	resp, err := saasApi.Post(
		os.Getenv("SAAS_HOST")+"/saas/api/v1/friend/getListFriend",
		nil,
		header,
		nil,
	)

	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	var resBody response.SaasResponse
	json.NewDecoder(resp.Body).Decode(&resBody)
	log.InfoLogger.Println(resBody)
	data := resBody.Data.([]interface{})
	return data, nil
}

func GetSaIdFromToken(token string) (string, error) {
	citizen, err := GetCitizen(token)
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", err
	}
	return citizen["saId"].(string), nil
}

func GetSaasAccessToken(username string, password string) (string, string, error) {
	resp, err := saasApi.Post(
		"https://keycloak.summonersarena.io/realms/summonersarena/protocol/openid-connect/token",
		ioutil.NopCloser(
			strings.NewReader(
				"username="+username+
					"&password="+password+
					"&grant_type=password&client_id="+os.Getenv("CLIENT_ID")+
					"&client_secret="+os.Getenv("CLIENT_SECRET")+"")),
		http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
		},
		nil,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", "", err
	}
	defer resp.Body.Close()
	var body map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", "", err
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", "", err
	}
	accessToken, ok := body["access_token"]
	if !ok {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		return "", "", errors.New(types.ErrorInvalidInput)
	}
	refreshToken, ok := body["refresh_token"]
	if !ok {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		return "", "", errors.New(types.ErrorInvalidInput)
	}
	return accessToken.(string), refreshToken.(string), nil

}
