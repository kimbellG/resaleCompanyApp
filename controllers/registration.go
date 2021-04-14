package controllers

import (
	"cw/models"
	"log"
	"net/http"
)

type Env struct {
	Reg Registrar
}

type Registrar interface {
	RegisterUser(userInfo *models.RegistrationInformation) error
}

//Обработать приход json с клиента
func (env *Env) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	result, err := models.DecodingAnswerForRegistration(r)
	if err != nil {
		log.Printf("registration: decoding json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("registering a new user: ", result.AuthInfo.Login)

	err = env.Reg.RegisterUser(result)
	if err != nil {
		log.Printf("registration: &v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens := result.AuthInfo.GenerateToken()
	nw, err := tokens.EncodingTokensInHttpBody(w)
	if err != nil {
		log.Fatalln("encoding tokens:", err)
	}
	w = nw
}
