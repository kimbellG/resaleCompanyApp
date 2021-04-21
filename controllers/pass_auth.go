package controllers

import (
	"cw/models"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Auth interface {
	GetUserInfo(user *models.PasswordAutheficationInfo) (*models.AuthorizationUserInformation, error)
}

func (env *Env) PasswordAuthentification(w http.ResponseWriter, r *http.Request) {
	PassAuthLog := logger.WithFields(log.Fields{
		"action": "password authentification",
	})

	logPass, err := models.DecodingPasswordAuthInfo(r)
	if err != nil {
		PassAuthLog.Debug("Incorrect request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	PassAuthLog.Debug(fmt.Sprintf("Connecting new user: %v", logPass.Login))

	userInfo, err := env.UserAuth.GetUserInfo(logPass)
	if err != nil {
		PassAuthLog.Debug(fmt.Sprintf("%v", err))
		http.Error(w, fmt.Sprintf("%v", err), http.StatusForbidden)
		return
	}

	tokens := userInfo.GenerateToken()
	w, err = tokens.EncodingTokensInHttpBody(w)
	if err != nil {
		log.Fatalln("encoding tokens:", err)
	}
}
