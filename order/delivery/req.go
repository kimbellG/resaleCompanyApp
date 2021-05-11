package delivery

import (
	"cw/order"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc order.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/order/add", h.Add)
	router.HandleFunc("/order/get", h.Gets)
	router.HandleFunc("/order/get/interval", h.GetInInterval)
	router.HandleFunc("/order/update/status", h.UpdateStatus)
	router.HandleFunc("/order/filter", h.Filter)
}
