package delivery

import (
	"cw/logger"
	"cw/product"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Handler struct {
	useCase product.UseCase
}

func NewHandler(cases product.UseCase) *Handler {
	return &Handler{
		useCase: cases,
	}
}

type AddRequests struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	addLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "add product record"},
	)

	newRecord := &AddRequests{}
	if err := decodingJson(r, newRecord); err != nil {
		addLogger.Debugf("Invalid json in request: %v", err)
		http.Error(w, "Invalid client json", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Add(r.Context(), reqToProduct(newRecord)); err != nil {
		addLogger.Debug(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func reqToProduct(req *AddRequests) *product.Product {
	return &product.Product{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
	}
}

func decodingJson(r *http.Request, strct interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(strct); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Gets(w http.ResponseWriter, r *http.Request) {
	getLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Get provider"},
	)

	resultPrv, err := h.useCase.Gets(r.Context())
	if err != nil {
		getLogger.Fatalln(err)
		os.Exit(1)
	}

	if err := encodingInBody(&w, resultPrv); err != nil {
		getLogger.Fatal(err)
		os.Exit(1)
	}
}

func encodingInBody(w *http.ResponseWriter, products []product.Product) error {
	result := new([]AddRequests)

	for _, val := range products {
		*result = append(*result, *provToClient(&val))
	}

	if err := encodingJson(*w, result); err != nil {
		return fmt.Errorf("encoding is failed: %v", err)
	}

	return nil
}

func provToClient(prov *product.Product) *AddRequests {
	return &AddRequests{
		Id:          prov.Id,
		Name:        prov.Name,
		Description: prov.Description,
	}
}

func encodingJson(w http.ResponseWriter, strct interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(strct); err != nil {
		return err
	}

	return nil
}

type Field struct {
	Key   string
	Value interface{}
}
type UpdateRequest struct {
	Id     int
	Fields *[]Field
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	updateLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Update provider"},
	)

	fields := new([]UpdateRequest)

	if err := decodingJson(r, fields); err != nil {
		updateLogger.Debugf("Invalid request: %v", err)
		http.Error(w, "Incorrect json update body", http.StatusBadRequest)
		return
	}

	for _, field := range *fields {
		if err := h.useCase.Update(r.Context(), field.Id, fieldToMap(*field.Fields)); err != nil {
			updateLogger.Debug(err)
			http.Error(w, "Invalid update request", http.StatusBadRequest)
			return
		}
	}
}

func fieldToMap(f []Field) map[string]interface{} {
	result := make(map[string]interface{})

	for _, v := range f {
		result[v.Key] = v.Value
	}

	return result
}

type deleteRequest struct {
	Code int
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	deleteLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Delete Provider"},
	)

	req := &deleteRequest{}
	if err := decodingJson(r, req); err != nil {
		deleteLogger.Debugf("Invalid request: %v", err)
		http.Error(w, "Incorrect update request", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Delete(r.Context(), req.Code); err != nil {
		deleteLogger.Debug(err)
		http.Error(w, "Invalid delete request", http.StatusBadRequest)
		return
	}
}

func (h *Handler) Filter(w http.ResponseWriter, r *http.Request) {
	filterLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "filter client"},
	)

	if err := r.ParseForm(); err != nil {
		filterLogger.Debugf("Incorrect form from query: %v", err)
		http.Error(w, "Incorrect URL-forms.", http.StatusBadRequest)
		return
	}

	if len(r.Form) != 1 {
		filterLogger.Debugf("Incorrect length of form")
		http.Error(w, "Incorrect URL forms", http.StatusBadRequest)
		return
	}

	for k, v := range r.Form {
		resultCl, err := h.useCase.Filter(r.Context(), k, v[0])
		if err != nil {
			filterLogger.Debugf("%v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := encodingInBody(&w, resultCl); err != nil {
			filterLogger.Fatalf("Invalid encoding: %v", err)
			os.Exit(1)
		}

	}
}
