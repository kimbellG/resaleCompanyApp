package controllers

import (
	"fmt"
	"net/http"
)

type Env struct {
	Reg      Registrar
	UserAuth Auth
}

func TestingToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "Your id is %v", user)
}
