package delivery

import (
	"cw/rang"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(router *mux.Router, uc rang.Usecase) {
	h := NewHandler(uc)

	router.HandleFunc("/rank/add/problem", h.PUTProblem)
	router.HandleFunc("/rank/add/mark", h.PUTMarks)
	router.HandleFunc("/rank/get", h.Gets)
	router.HandleFunc("/rank/get_by_id", h.GetProblem)
}
