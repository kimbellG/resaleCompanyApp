package delivery

import (
	"cw/auth"
	"cw/logger"
	"encoding/json"
	"net/http"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(cases auth.UseCase) *Handler {
	return &Handler{
		useCase: cases,
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
	Token string
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

	if err := h.useCase.SignUp(r.Context(), newUser.Login, newUser.Password, newUser.Name); err != nil {
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

	token, err := h.useCase.SignIn(r.Context(), userSignIn.Login, userSignIn.Password)
	if err != nil {
		inLogger.Debugf("sign in is failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SendToken(w, token)
}

func decodingAnswerForSignIn(r *http.Request) (*SignInInput, error) {
	result := &SignInInput{}

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func SendToken(w http.ResponseWriter, token string) {
	tokenResult := &TokenType{Token: token}

	if err := json.NewEncoder(w).Encode(tokenResult); err != nil {
		panic("Incorrect token encoding")
	}
}
