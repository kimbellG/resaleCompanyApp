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

type AddRequests struct {
	VendorCode     int
	Name           string
	UNP            string `json:"unp"`
	TermsOfPayment string `json:"terms_of_payment"`
	Address        string
	PhoneNumber    string `json:"phone_number"`
	Email          string
	WebSite        string `json:"web_site"`
}

func (h *Handler) AddProvider(w http.ResponseWriter, r *http.Request) {
	addLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "add provider record"},
	)

	newRecord := &AddRequests{}
	if err := decodingJson(r, newRecord); err != nil {
		addLogger.Debugf("Invalid json in request: %v", err)
		http.Error(w, "Invalid provider json", http.StatusBadRequest)
		return
	}

	if err := h.useCase.AddProvider(r.Context(), reqToProv(newRecord)); err != nil {
		addLogger.Debug(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func reqToProv(req *AddRequests) *provider.Provider {
	return &provider.Provider{
		VendorCode:     req.VendorCode,
		Name:           req.Name,
		UNP:            req.UNP,
		TermsOfPayment: req.TermsOfPayment,
		PhoneNumber:    req.PhoneNumber,
		Address:        req.Address,
		Email:          req.Email,
		WebSite:        req.WebSite,
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

	resultPrv, err := h.useCase.GetProviders(r.Context())
	if err != nil {
		getLogger.Fatalln(err)
		os.Exit(1)
	}

	result := new([]AddRequests)

	for _, val := range resultPrv {
		*result = append(*result, *provToReq(&val))
	}

	if err := encodingJson(w, result); err != nil {
		getLogger.Fatalf("Encoding is failed: %v", err)
		os.Exit(1)
	}
}

func provToReq(prov *provider.Provider) *AddRequests {
	return &AddRequests{
		VendorCode:     prov.VendorCode,
		Name:           prov.Name,
		UNP:            prov.UNP,
		TermsOfPayment: prov.TermsOfPayment,
		PhoneNumber:    prov.PhoneNumber,
		Address:        prov.Address,
		Email:          prov.Email,
		WebSite:        prov.WebSite,
	}
}

func encodingJson(w http.ResponseWriter, strct interface{}) error {
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
	Code   int
	Fields *[]Field
}

func (h *Handler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
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
		if err := h.useCase.UpdateProvider(r.Context(), field.Code, fieldToMap(*field.Fields)); err != nil {
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

	if err := h.useCase.DeleteProvider(r.Context(), req.Code); err != nil {
		deleteLogger.Debug(err)
		http.Error(w, "Invalid delete request", http.StatusBadRequest)
		return
	}
}
