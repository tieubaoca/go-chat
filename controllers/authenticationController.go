package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
)

func Authentication(c *gin.Context) {
	tokenString := utils.GetAccessTokenByReq(c.Request)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorTokenEmpty,
			Data:    "",
		})
		return
	}
	token, err := utils.ParseUnverified(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data:    token.Claims,
	})
}

func GetAccessToken(c *gin.Context) {
	var getAccTokenReq request.GetAccessTokenReq
	err := json.NewDecoder(c.Request.Body).Decode(&getAccTokenReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	if getAccTokenReq.Username == "" || getAccTokenReq.Password == "" {
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	accessToken, refreshToken, err := services.GetAccessToken(getAccTokenReq.Username, getAccTokenReq.Password)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	c.SetCookie(
		"access-token",
		accessToken,
		10*60,
		"/",
		os.Getenv("DOMAIN"),
		false,
		false,
	)
	c.SetCookie(
		"refresh-token",
		refreshToken,
		10*60,
		"/",
		os.Getenv("DOMAIN"),
		false,
		false,
	)
	_, err = utils.Parse(accessToken)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data: map[string]string{
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		},
	})

}

// func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
// 	accessTokenCookie, err := r.Cookie("access-token")
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	accessToken := accessTokenCookie.Value
// 	if accessToken == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, types.StatusError, nil, "Access token is empty")
// 		return
// 	}
// 	refreshedAccessToken, err := services.RefreshAccessToken(accessToken)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "access_token",
// 		Value:   refreshedAccessToken,
// 		Expires: time.Now().Add(24 * time.Hour),
// 	})
// 	response.Res(w, types.StatusSuccess, nil, "Refresh access token successfully")
// }
