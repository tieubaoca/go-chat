package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils/httpHelper"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

type SaasService interface {
	GetSaasAccessToken(username string, password string) (string, string, error)
	GetCitizen(token string) (map[string]interface{}, error)
	GetListFriendInfo(saasAccessToken string, paginationReq request.PaginationOnlineFriendReq) ([]interface{}, error)
	GetAllFriends(token string) ([]interface{}, error)
	GetSaId(token string) (string, error)
}

type saasService struct{}

func NewSaasService() *saasService {
	return &saasService{}
}

func (s *saasService) GetSaasAccessToken(username string, password string) (string, string, error) {
	resp, err := httpHelper.Post(
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

func (s *saasService) GetCitizen(token string) (map[string]interface{}, error) {
	header := http.Header{}
	header.Add(
		"Authorization",
		"Bearer "+token,
	)
	res, err := httpHelper.Get(
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
	return data, nil
}

func (s *saasService) GetListFriendInfo(saasAccessToken string, paginationReq request.PaginationOnlineFriendReq) ([]interface{}, error) {
	header := http.Header{}
	header.Add(
		"Authorization",
		"Bearer "+saasAccessToken,
	)
	query := url.Values{
		"page": []string{fmt.Sprint(paginationReq.Page)},
		"size": []string{fmt.Sprint(paginationReq.Size)},
	}

	resp, err := httpHelper.Get(
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

func (s *saasService) GetAllFriends(token string) ([]interface{}, error) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)
	resp, err := httpHelper.Post(
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

func (s *saasService) GetSaId(token string) (string, error) {
	citizen, err := s.GetCitizen(token)
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", err
	}
	return citizen["saId"].(string), nil
}
