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
