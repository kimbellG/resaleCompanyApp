package delivery

import (
	"cw/product"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc product.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/product/add", h.Add)
	router.HandleFunc("/product/get", h.Gets)
	router.HandleFunc("/product/delete", h.Delete)
	router.HandleFunc("/product/update", h.Update)
	router.HandleFunc("/product/filter", h.Filter)
}
