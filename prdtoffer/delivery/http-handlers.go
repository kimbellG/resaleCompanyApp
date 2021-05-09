package delivery

import (
	"cw/logger"
	"cw/prdtoffer"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Handler struct {
	useCase prdtoffer.UseCase
}

func NewHandler(cases prdtoffer.UseCase) *Handler {
	return &Handler{
		useCase: cases,
	}
}

type AddRequests struct {
	Provider string  `json:"provider"`
	Product  string  `json:"product"`
	Cost     float32 `json:"cost"`
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	addLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "add offer record"},
	)

	newRecord := &AddRequests{}
	if err := decodingJson(r, newRecord); err != nil {
		addLogger.Debugf("Invalid json in request: %v", err)
		http.Error(w, "Invalid offer json", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Add(r.Context(), reqToProduct(newRecord)); err != nil {
		addLogger.Debug(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func reqToProduct(req *AddRequests) *prdtoffer.Offer {
	return &prdtoffer.Offer{
		ProductName:  req.Product,
		ProviderName: req.Provider,
		Cost:         req.Cost,
	}
}

func decodingJson(r *http.Request, strct interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(strct); err != nil {
		return err
	}

	return nil
}

type ProductInfo struct {
	Product string
}

func (h *Handler) Gets(w http.ResponseWriter, r *http.Request) {
	getLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Get offer"},
	)

	if err := r.ParseForm(); err != nil {
		getLogger.Debugf("Incorrect form from query: %v", err)
		http.Error(w, "Incorrect URL-forms.", http.StatusBadRequest)
		return
	}

	product, isProduct := r.Form["product"]
	provider, isProvider := r.Form["provider"]

	resultPrv, err := *new([]prdtoffer.Offer), error(nil)
	if isProduct {
		resultPrv, err = h.useCase.GetOfferForProduct(r.Context(), product[0])
		if err != nil {
			getLogger.Debugf("for product: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if isProvider {
		resultPrv, err = h.useCase.GetOffersOfProvider(r.Context(), provider[0])

		if err != nil {
			getLogger.Debugf("of provider: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := encodingInBody(&w, resultPrv); err != nil {
		getLogger.Fatal(err)
		os.Exit(1)
	}
}

func encodingInBody(w *http.ResponseWriter, products []prdtoffer.Offer) error {
	result := new([]AddRequests)

	for _, val := range products {
		*result = append(*result, *provToClient(&val))
	}

	if err := encodingJson(*w, result); err != nil {
		return fmt.Errorf("encoding is failed: %v", err)
	}

	return nil
}

func provToClient(offer *prdtoffer.Offer) *AddRequests {
	return &AddRequests{
		Product:  offer.ProductName,
		Provider: offer.ProviderName,
		Cost:     offer.Cost,
	}
}

func encodingJson(w http.ResponseWriter, strct interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(strct); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	updateLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Update offer cost"},
	)

	fields := &AddRequests{}

	if err := decodingJson(r, fields); err != nil {
		updateLogger.Debugf("Invalid request: %v", err)
		http.Error(w, "Incorrect json update body", http.StatusBadRequest)
		return
	}

	if err := h.useCase.UpdateCost(r.Context(), fields.Product, fields.Provider, fields.Cost); err != nil {
		updateLogger.Debug(err)
		http.Error(w, "Invalid update request", http.StatusBadRequest)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	deleteLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "Delete offer"},
	)

	req := &AddRequests{}
	if err := decodingJson(r, req); err != nil {
		deleteLogger.Debugf("Invalid request: %v", err)
		http.Error(w, "Incorrect update request", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Delete(r.Context(), req.Product, req.Provider); err != nil {
		deleteLogger.Debug(err)
		http.Error(w, "Invalid delete request", http.StatusBadRequest)
		return
	}
}
