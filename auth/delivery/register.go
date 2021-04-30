package delivery

import (
	"cw/auth"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc auth.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/sign-up", h.SignUp)
}
