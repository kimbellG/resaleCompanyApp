package controllers

import (
	neasted "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.New()
	logger.SetFormatter(&neasted.Formatter{
		HideKeys: true,
	})
	logger.SetLevel(log.DebugLevel)
}
