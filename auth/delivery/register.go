package delivery

import (
	"cw/auth"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc auth.UseCase, auc auth.AdminUseCase) {
	h := NewHandler(uc, auc)

	router.HandleFunc("/sign-up", h.SignUp)
	router.HandleFunc("/sign-in", h.SignIn)

	router.HandleFunc("/user-control/confirm", h.ConfirmUser)
	router.HandleFunc("/user-control/disable", h.DisableUser)
	router.HandleFunc("/user-control/get/all", h.GetAllUsers)
	router.HandleFunc("/user-control/get/by/login", h.GetUser)
	router.HandleFunc("/user-control/update", h.UpdateUser)
	router.HandleFunc("/user-control/delete", h.DeleteUser)
}
