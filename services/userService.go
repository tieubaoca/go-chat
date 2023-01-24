package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PaginationOnlineFriends(saId string, paginationReq request.PaginationOnlineFriendReq) ([]interface{}, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()

	body := map[string]string{"saId": saId}

	jsonBody := new(bytes.Buffer)
	json.NewEncoder(jsonBody).Encode(body)

	req := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   os.Getenv("SAAS_HOST"),
			Path:   "/saas/api/v1/friend/getListFriendInfo",
			RawQuery: url.Values{
				"page": []string{fmt.Sprint(paginationReq.Page)},
				"size": []string{fmt.Sprint(paginationReq.Size)},
			}.Encode(),
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},

		Body: io.NopCloser(jsonBody),
	}
	client := &http.Client{}

	resp, err := client.Do(&req)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	var resBody response.SaasResponse
	json.NewDecoder(resp.Body).Decode(&resBody)
	log.InfoLogger.Println(resBody)
	data := resBody.Data.([]interface{})
	saIds := make([]string, 0)
	for _, v := range data {
		saIds = append(saIds, v.(map[string]interface{})["saId"].(string))
	}
	userStatus := FindUserStatusInSaIdList(saIds)
	for _, v := range data {
		friendSaId := v.(map[string]interface{})["saId"].(string)
		if friendStatus, ok := userStatus[friendSaId]; ok {
			v.(map[string]interface{})["online"] = friendStatus.IsActive
			v.(map[string]interface{})["lastSeen"] = friendStatus.LastSeen
		} else {
			v.(map[string]interface{})["online"] = false
			v.(map[string]interface{})["lastSeen"] = 0
		}
	}
	return data, nil
}

func FindUserStatusInSaIdList(saIds []string) map[string]models.UserOnlineStatus {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.UserOnlineStatusCollection)
	result, err := coll.Find(context.TODO(), bson.D{{"saId", bson.D{{"$in", saIds}}}})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil
	}
	var users []models.UserOnlineStatus
	if err = result.All(context.TODO(), &users); err != nil {
		log.ErrorLogger.Println(err)
		return nil
	}
	mapUser := make(map[string]models.UserOnlineStatus)
	for _, v := range users {
		mapUser[v.SaId] = v
	}
	return mapUser

}

func FindUserStatusBySaId(saId string) models.UserOnlineStatus {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.UserOnlineStatusCollection)
	result := coll.FindOne(context.TODO(), bson.D{{"saId", saId}})
	var user models.UserOnlineStatus
	if err := result.Decode(&user); err != nil {
		log.ErrorLogger.Println(err)
	}
	return user
}

func UpdateUserStatus(saId string, isActive bool, lastSeen primitive.DateTime) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.UserOnlineStatusCollection)
	if FindUserStatusBySaId(saId).SaId == "" {
		_, err := coll.InsertOne(context.TODO(), bson.M{
			"saId":     saId,
			"isActive": isActive,
			"lastSeen": lastSeen,
		})
		if err != nil {
			log.ErrorLogger.Println(err)
		}
		return
	}
	_, err := coll.UpdateOne(context.TODO(), bson.D{{"saId", saId}}, bson.D{
		{"$set", bson.D{
			{"isActive", isActive},
			{"lastSeen", lastSeen},
		}},
	})
	if err != nil {
		log.ErrorLogger.Println(err)
	}

}
