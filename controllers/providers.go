package controllers

import (
	"cw/models"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type Provider interface {
	InsertProviderInDB(pr *models.Provider) error
}

func (env *Env) CreateProviderController(w http.ResponseWriter, r *http.Request) {
	ProviderLogger := logger.WithFields(log.Fields{"action": "create provider"})

	prov, err := models.DecodingProvider(r)
	if err != nil {
		ProviderLogger.Debug("parsing json: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ProviderLogger.Debugf("Creating provider in db: name=%v", prov.Name)

	if err := env.Prov.InsertProviderInDB(prov); err != nil {
		ProviderLogger.Errorf("Invalid querty to db: %v", err)
		os.Exit(1)
	}

	ProviderLogger.Debugf("CreatingComplete")
}
