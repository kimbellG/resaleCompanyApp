package delivery

import (
	"cw/logger"
	"cw/models"
	"cw/order"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Handler struct {
	usecase order.UseCase
}

func NewHandler(cases order.UseCase) *Handler {
	return &Handler{
		usecase: cases,
	}
}

type AddRequests struct {
	Id           int    `json:"id"`
	Offers       []int  `json:"offers"`
	ClientId     int    `json:"client_id"`
	ManagerLogin string `json:"manager_login"`
	OrderDate    string `json:"order_date"`
	Quantity     int    `json:"quantity"`
	Status       string `json:"status"`
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	addLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "add order information"},
	)

	ordr := &AddRequests{}
	if err := decodingJson(r, ordr); err != nil {
		addLogger.Debugf("decoding json: %v", err)
		http.Error(w, "Incorrect json request", http.StatusBadRequest)
		return
	}

	newOrder, err := reqToMod(ordr)
	if err != nil {
		addLogger.Debugf("request to order model: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.Add(r.Context(), newOrder); err != nil {
		addLogger.Debugf("usecase: %v", err)
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

func reqToMod(req *AddRequests) (*models.Order, error) {
	orderDate, err := time.Parse(time.UnixDate, req.OrderDate)
	if err != nil {
		return nil, fmt.Errorf("parsing time: %v", err)
	}
	return &models.Order{
		Id:           req.Id,
		Offers:       req.Offers,
		ClientId:     req.ClientId,
		ManagerLogin: req.ManagerLogin,
		OrderDate:    orderDate,
		Quantity:     req.Quantity,
		Status:       req.Status,
	}, nil
}

type Output struct {
	Id int `json:"id"`
}

func (h *Handler) Gets(w http.ResponseWriter, r *http.Request) {

}
