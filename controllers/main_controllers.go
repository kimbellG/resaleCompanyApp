package controllers

import (
	"fmt"
	"net/http"
)

func TestingToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "Your id is %v", user)
}
