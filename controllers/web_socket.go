package controllers

import (
	"net/http"

	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	username := session.Values["username"].(string)
	services.HandleWebSocket(w, r, username, session.ID)
}
