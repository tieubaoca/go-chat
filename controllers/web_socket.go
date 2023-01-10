package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	log.Info(token)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	services.HandleWebSocket(w, r, utils.GetSaIdFromToken(token), utils.GetSessionIdFromToken(token))
}
