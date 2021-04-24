package controllers

import (
	"cw/models"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (env *Env) CreateProviderController(w http.ResponseWriter, r *http.Request) {
	ProviderLogger := logger.WithFields(log.Fields{"action": "create provider"})

	prov, err := models.DecodingProvider(r)
	if err != nil {
		ProviderLogger.Debug("parsing json: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

}
