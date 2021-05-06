package delivery

import (
	"cw/logger"
	"cw/provider"
	"encoding/json"
	"net/http"
	"os"
)

type Handler struct {
	useCase provider.UseCase
}

func NewHandler(cases provider.UseCase) *Handler {
	return &Handler{
		useCase: cases,
	}
}

func (h *Handler) AddProvider(w http.ResponseWriter, r *http.Request) {
	addLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "add provider record"},
	)

	newRecord := &provider.Provider{}
	if err := decodingJson(r, newRecord); err != nil {
		addLogger.Debugf("Invalid json in request: %v", err)
		http.Error(w, "Invalid provider json", http.StatusBadRequest)
		return
	}

	if err := h.useCase.AddProvider(r.Context(), newRecord); err != nil {
		addLogger.Debug(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func decodingJson(r *http.Request, strct interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(strct); err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetProviders(w http.ResponseWriter, r *http.Request) {
	getLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Get provider"},
	)

	result, err := h.useCase.GetProviders(r.Context())
	if err != nil {
		getLogger.Fatalln(err)
		os.Exit(1)
	}

	if err := encodingJson(w, result); err != nil {
		getLogger.Fatalf("Encoding is failed: %v", err)
		os.Exit(1)
	}
}

func encodingJson(w http.ResponseWriter, strct interface{}) error {
	if err := json.NewEncoder(w).Encode(strct); err != nil {
		return err
	}

	return nil
}

type Field struct {
	key   string
	value interface{}
}
type UpdateRequest struct {
	code   int
	fields []Field
}

func (h *Handler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	updateLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Update provider"},
	)

	fields := new([]UpdateRequest)

	if err := decodingJson(r, *fields); err != nil {
		updateLogger.Debugf("Invalid request: %v", err)
		http.Error(w, "Incorrect json update body", http.StatusBadRequest)
		return
	}

	for _, field := range *fields {
		if err := h.useCase.UpdateProvider(r.Context(), field.code, fieldToMap(field.fields)); err != nil {
			updateLogger.Debug(err)
			http.Error(w, "Invalid update request", http.StatusBadRequest)
			return
		}
	}
}

func fieldToMap(f []Field) map[string]interface{} {
	result := *new(map[string]interface{})

	for _, v := range f {
		result[v.key] = v.value
	}

	return result
}

type deleteRequest struct {
	code int
}

func (h *Handler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	deleteLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Delete Provider"},
	)

	req := &deleteRequest{}
	if err := decodingJson(r, req); err != nil {
		deleteLogger.Debugf("Invalid request: %v", err)
		http.Error(w, "Incorrect update request", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteProvider(r.Context(), req.code); err != nil {
		deleteLogger.Debug(err)
		http.Error(w, "Invalid delete request", http.StatusBadRequest)
		return
	}
}
