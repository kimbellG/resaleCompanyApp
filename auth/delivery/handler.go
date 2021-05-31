package delivery

import (
	"cw/auth"
	"cw/logger"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	UserUseCase  auth.UseCase
	AdminUseCase auth.AdminUseCase
	logg         *logrus.Entry
}

func NewHandler(userCases auth.UseCase, adminCases auth.AdminUseCase) *Handler {
	return &Handler{
		UserUseCase:  userCases,
		AdminUseCase: adminCases,
		logg: logger.NewLoggerWithFields(
			map[string]interface{}{"action": "user control"},
		),
	}
}

type SignInput struct {
	Login    string
	Password string
	Name     string
}

type SignInInput struct {
	Login    string
	Password string
}

type TokenType struct {
	Token  string `json:"token"`
	Access string `json:"access"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	upLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "sign up"})

	newUser, err := decodingAnswerForRegistration(r)

	if err != nil {
		http.Error(w, "Invalid json information for sign up", http.StatusBadRequest)
		upLogger.Errorf("invalid json decoding: %v", err)
		return
	}

	if err := h.UserUseCase.SignUp(r.Context(), newUser.Login, newUser.Password, newUser.Name); err != nil {
		upLogger.Errorf("invalid sign up in usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func decodingAnswerForRegistration(r *http.Request) (*SignInput, error) {
	result := &SignInput{}

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	inLogger := logger.NewLoggerWithFields(map[string]interface{}{"action": "Sign in"})

	userSignIn, err := decodingAnswerForSignIn(r)

	if err != nil {
		inLogger.Debugf("Incorrect json message: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.UserUseCase.SignIn(r.Context(), userSignIn.Login, userSignIn.Password)
	if err != nil {
		inLogger.Debugf("sign in is failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SendToken(w, token.Token, token.Access)
}

func decodingAnswerForSignIn(r *http.Request) (*SignInInput, error) {
	result := &SignInInput{}

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func SendToken(w http.ResponseWriter, token, access string) {
	tokenResult := &TokenType{Token: token, Access: access}

	if err := json.NewEncoder(w).Encode(tokenResult); err != nil {
		panic("Incorrect token encoding")
	}
}
