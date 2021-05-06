package delivery

import (
	"cw/provider"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc provider.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/provider/add", h.AddProvider)
	router.HandleFunc("/provider/get", h.GetProviders)
	router.HandleFunc("/provider/delete", h.DeleteProvider)
	router.HandleFunc("/provider/update", h.UpdateProvider)
}
