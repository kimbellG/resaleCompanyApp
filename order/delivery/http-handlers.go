package delivery

import (
	"cw/client"
	"cw/logger"
	"cw/models"
	"cw/order"
	"cw/prdtoffer"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase     order.UseCase
	OrderLogger *logrus.Entry
}

func NewHandler(cases order.UseCase) *Handler {
	return &Handler{
		usecase:     cases,
		OrderLogger: logger.NewLoggerWithFields(map[string]interface{}{"action": "order"}),
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
	Id        int               `json:"id"`
	Offers    []prdtoffer.Offer `json:"offers"`
	Client    client.Client     `json:"client"`
	Manager   string            `json:"manager"`
	OrderDate string            `json:"order_date"`
	Quantity  int               `json:"quantity"`
	Status    string            `json:"status"`
}

func (h *Handler) Gets(w http.ResponseWriter, r *http.Request) {

	result, err := h.usecase.Gets(r.Context())
	h.get(w, result, err)
}

func (h *Handler) get(w http.ResponseWriter, orders []order.OrderOutput, err error) {
	if err != nil {
		h.OrderLogger.Debugf("get: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if err := encodintJson(w, arrOrderToOutput(orders)); err != nil {
		h.OrderLogger.Debugf("get: encoding in body: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func arrOrderToOutput(orders []order.OrderOutput) []Output {
	result := *new([]Output)
	for _, order := range orders {
		result = append(result, *OrderToOutput(&order))
	}

	return result
}

func OrderToOutput(order *order.OrderOutput) *Output {
	return &Output{
		Id:        order.Id,
		Offers:    order.Offers,
		Client:    order.Client,
		Manager:   order.Manager,
		OrderDate: order.OrderDate.Format(time.UnixDate),
		Quantity:  order.Quantity,
		Status:    order.Status,
	}
}

func encodintJson(w http.ResponseWriter, value interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return err
	}

	return nil
}

type GetTimeRequest struct {
	Start string
	End   string
}

func (h *Handler) GetInInterval(w http.ResponseWriter, r *http.Request) {
	dateReq := &GetTimeRequest{}

	if err := decodingJson(r, dateReq); err != nil {
		h.OrderLogger.Debugf("get in interval: decoding json body: %v", err)
		http.Error(w, "Incorrect json request", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.UnixDate, dateReq.Start)
	if err != nil {
		h.OrderLogger.Debugf("get in interval: parsing start date: %v", err)
		http.Error(w, "incorrect start format date", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.UnixDate, dateReq.End)
	if err != nil {
		h.OrderLogger.Debugf("get in interval: parsing end date: %v", err)
		http.Error(w, "incorrect end format date", http.StatusBadRequest)
		return
	}

	result, err := h.usecase.GetInInterval(r.Context(), startTime, endTime)
	h.get(w, result, err)
}

type UpdateRequest struct {
	Id     int
	Status string
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	req := UpdateRequest{}

	if err := decodingJson(r, req); err != nil {
		h.OrderLogger.Debugf("update status: invalid json body: %v", err)
		http.Error(w, "incorrect json body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.UpdateStatus(r.Context(), req.Id, req.Status); err != nil {
		h.OrderLogger.Debugf("update status: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) Filter(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.OrderLogger.Debugf("filter: parse form: %v", err)
		http.Error(w, "incorrect url form", http.StatusBadRequest)
		return
	}

	if len(r.Form) != 1 {
		h.OrderLogger.Debugf("filter: incorrect quantity of form")
		http.Error(w, "incorrect quantity of form", http.StatusBadRequest)
	}

	result, filterErr := *new([]order.OrderOutput), error(nil)
	for key, value := range r.Form {
		switch key {
		case "id", "clientId", "manager":
			val_int, err := strconv.Atoi(value[0])
			if err != nil {
				h.OrderLogger.Debugf("filter: atoi %v: %v", key, err)
				http.Error(w, fmt.Sprintf("incorrect %v value", key), http.StatusBadRequest)
				return
			}
			result, filterErr = h.usecase.Filter(r.Context(), key, val_int)
		case "status":
			result, filterErr = h.usecase.Filter(r.Context(), key, value)
		default:
			h.OrderLogger.Debug("filter: incorrect key value")
			http.Error(w, fmt.Sprintf("key(%v) not exists", key), http.StatusBadRequest)
			return
		}
	}

	h.get(w, result, filterErr)
}
