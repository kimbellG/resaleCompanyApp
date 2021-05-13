package delivery

import (
	"cw/logger"
	"cw/models"
	"cw/rang"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	useCase    rang.Usecase
	RankLogger *logrus.Entry
}

func NewHandler(cases rang.Usecase) *Handler {
	return &Handler{
		useCase:    cases,
		RankLogger: logger.NewLoggerWithFields(map[string]interface{}{"action": "Rank method"}),
	}
}

func (h *Handler) PUTProblem(w http.ResponseWriter, r *http.Request) {
	reqInput := &models.ProblemInput{}
	if err := decodingJSONFromRequest(r, reqInput); err != nil {
		h.RankLogger.Debugf("put problem: decoding json: %v", err)
		http.Error(w, "incorrect request json", http.StatusBadRequest)
		return
	}

	if err := h.useCase.Add(r.Context(), reqInput); err != nil {
		h.RankLogger.Debugf("put problem: add usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func decodingJSONFromRequest(r *http.Request, result interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

type MarkRequest struct {
	ProblemID     int     `json:"problem_id"`
	AlternativeID int     `json:"alternative_id"`
	ExpertLogin   string  `json:"expert_login"`
	Mark          float32 `json:"mark"`
}

func (h *Handler) PUTMarks(w http.ResponseWriter, r *http.Request) {
	req := &MarkRequest{}
	if err := decodingJSONFromRequest(r, req); err != nil {
		h.RankLogger.Debugf("put marks: decoding json: %v", err)
		http.Error(w, "incorrcet JSON request", http.StatusBadRequest)
		return
	}

	if err := h.useCase.AddAlternativeMark(r.Context(), markReqToModel(req)); err != nil {
		h.RankLogger.Debugf("put mark: usecase add mark: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func markReqToModel(req *MarkRequest) *models.AlternativeMarkInput {
	return &models.AlternativeMarkInput{
		ProblemId:     req.ProblemID,
		AlternativeId: req.AlternativeID,
		ExpertLogin:   req.ExpertLogin,
		Mark:          req.Mark,
	}
}

func (h *Handler) Gets(w http.ResponseWriter, r *http.Request) {
	result, err := h.useCase.Gets(r.Context())
	if err != nil {
		h.RankLogger.Debugf("get all problem: gets usecase: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := encodingJSONInResponse(w, result); err != nil {
		h.RankLogger.Debugf("get all problem: encoding json: %v", err)
		http.Error(w, "trouble with encoding json", http.StatusServiceUnavailable)
		return
	}
}

func encodingJSONInResponse(w http.ResponseWriter, item interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetProblem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.RankLogger.Debugf("get problem: parse url form: %v", err)
		http.Error(w, "incorrect url get form", http.StatusBadRequest)
		return
	}

	idStr, ok := r.Form["id"]
	if !ok {
		h.RankLogger.Debugf("get problem: forms: not found id form")
		http.Error(w, "not fount id form", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		h.RankLogger.Debugf("get problem: atoi: %v", err)
		http.Error(w, "incorrect id value", http.StatusBadRequest)
		return
	}

	result, err := h.useCase.GetProblemReport(r.Context(), id)
	if err != nil {
		h.RankLogger.Debugf("get problem: usecase: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := encodingJSONInResponse(w, result); err != nil {
		h.RankLogger.Debugf("get problem: encoding json: %v", err)
		http.Error(w, "encoding has been failed", http.StatusServiceUnavailable)
		return
	}
}
