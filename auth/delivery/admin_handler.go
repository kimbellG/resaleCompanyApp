package delivery

import (
	"encoding/json"
	"net/http"
)

type ConfirmRequest struct {
	Username      string `json:"username"`
	AccessProfile string `json:"access_profile"`
}

func (h *Handler) ConfirmUser(w http.ResponseWriter, r *http.Request) {
	req := &ConfirmRequest{}

	if err := decodingJson(r, req); err != nil {
		h.logg.Debug("confrim user: incorrect json request: %v", err)
		http.Error(w, "incorrect json request", http.StatusBadRequest)
		return
	}

	if err := h.AdminUseCase.ConfirmUser(r.Context(), req.Username, req.AccessProfile); err != nil {
		h.logg.Debugf("confrim user: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func decodingJson(r *http.Request, result interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

type DefaultRequest struct {
	Username string `json:"username"`
}

func (h *Handler) DisableUser(w http.ResponseWriter, r *http.Request) {
	req := &DefaultRequest{}

	if err := decodingJson(r, req); err != nil {
		h.logg.Debug("disable user: incorrect json request: %v", err)
		http.Error(w, "incorrect json request", http.StatusBadRequest)
		return
	}

	if err := h.AdminUseCase.DisableUser(r.Context(), req.Username); err != nil {
		h.logg.Debug("disable user: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	result, err := h.AdminUseCase.GetAllUsers(r.Context())
	if err != nil {
		h.logg.Debug("get all users: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := encodintJson(w, result); err != nil {
		h.logg.Debugf("get all users: encoding in body: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func encodintJson(w http.ResponseWriter, value interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	req := &DefaultRequest{}

	if err := decodingJson(r, req); err != nil {
		h.logg.Debug("get user: incorrect json request: %v", err)
		http.Error(w, "incorrect json request", http.StatusBadRequest)
		return
	}

	result, err := h.AdminUseCase.GetUser(r.Context(), req.Username)
	if err != nil {
		h.logg.Debug("get user: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := encodintJson(w, result); err != nil {
		h.logg.Debugf("get user: encoding in body: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

type UpdateRequest struct {
	Username string `json:"username"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	req := &UpdateRequest{}

	if err := decodingJson(r, req); err != nil {
		h.logg.Debug("update user: incorrect json request: %v", err)
		http.Error(w, "incorrect json request", http.StatusBadRequest)
		return
	}

	if err := h.AdminUseCase.UpdateUser(r.Context(), req.Username, req.Key, req.Value); err != nil {
		h.logg.Debug("update user: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	req := &DefaultRequest{}

	if err := decodingJson(r, req); err != nil {
		h.logg.Debug("delete user: incorrect json request: %v", err)
		http.Error(w, "incorrect json request", http.StatusBadRequest)
		return
	}

	if err := h.AdminUseCase.DeleteUser(r.Context(), req.Username); err != nil {
		h.logg.Debug("delete user: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
