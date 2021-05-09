package delivery

import (
	"cw/prdtoffer"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc prdtoffer.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/offer/add", h.Add)
	router.HandleFunc("/offer/get", h.Gets)
	router.HandleFunc("/offer/delete", h.Delete)
	router.HandleFunc("/offer/update", h.Update)
}
