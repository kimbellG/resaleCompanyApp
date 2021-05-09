package delivery

import (
	"cw/client"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc client.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/client/add", h.AddClient)
	router.HandleFunc("/client/get", h.GetClients)
	router.HandleFunc("/client/delete", h.DeleteClient)
	router.HandleFunc("/client/update", h.UpdateClient)
	router.HandleFunc("/client/filter", h.FilterClient)
}
