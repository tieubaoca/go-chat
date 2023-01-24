package controllers

import (
	"net/http"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	log.InfoLogger.Println(token)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	services.HandleWebSocket(w, r, utils.GetSaIdFromToken(token), utils.GetSessionIdFromToken(token))
}
