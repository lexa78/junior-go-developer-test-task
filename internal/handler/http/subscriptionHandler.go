package http

import (
	"encoding/json"
	"net/http"
	"testTask/internal/domain"
	"testTask/internal/dto"
	"testTask/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const DateFormatFromRequest = "01-2006"

// парсит дату в формате "MM-YYYY"
func parseMonthYear(dateStr string) (time.Time, error) {
	return time.Parse(DateFormatFromRequest, dateStr)
}

// парсит *string в *time.Time
func parseMonthYearPtr(dateStr *string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}
	t, err := time.Parse(DateFormatFromRequest, *dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// парсит UUID из строки
func parseUUID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}

// безопасная запись JSON с обработкой ошибки
func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewHandler(svc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: svc}
}

func (h *SubscriptionHandler) RegisterRoutes(r chi.Router) {
	r.Post("/subscriptions", h.Create)
	r.Get("/subscriptions/{id}", h.Get)
	r.Patch("/subscriptions/{id}", h.Update)
	r.Delete("/subscriptions/{id}", h.Delete)
	r.Get("/subscriptions/list", h.List)
	r.Get("/subscriptions/total", h.Total)
}

// Create godoc
// @Summary Запись новой подписки
// @Description Создание записи о новой подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Тело запроса"
// @Success 201 {object} domain.Subscription
// @Failure 400 {string} string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start, err := time.Parse(DateFormatFromRequest, req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	var end *time.Time
	if req.EndDate != nil {
		t, err := time.Parse(DateFormatFromRequest, *req.EndDate)
		if err != nil {
			http.Error(w, "invalid end_date format", http.StatusBadRequest)
			return
		}
		end = &t
	}

	userUuid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	sub := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userUuid,
		StartDate:   start,
		EndDate:     end,
	}

	if err := h.service.Create(r.Context(), sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sub, http.StatusCreated)
}

// Total godoc
// @Summary Подсчет суммарной стоимости всех подписок
// @Description Подсчет суммарной стоимости всех подписок с фильтрами по user и service
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Наименование сервиса"
// @Param from query string true "Начало периода. Формат MM-YYYY"
// @Param to query string true "Окончание периода. Формат MM-YYYY"
// @Success 200 {object} map[string]int
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) Total(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	serviceNameStr := r.URL.Query().Get("service_name")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		http.Error(w, "dates 'from' and 'to' are required", http.StatusBadRequest)
		return
	}

	from, err := parseMonthYear(fromStr)
	if err != nil {
		http.Error(w, "invalid date 'from'", http.StatusBadRequest)
		return
	}

	to, err := parseMonthYear(toStr)
	if err != nil {
		http.Error(w, "invalid date 'to'", http.StatusBadRequest)
		return
	}

	var userID *uuid.UUID
	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		userID = &uid
	}

	var serviceName *string
	if serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	filter := domain.TotalFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          to,
	}

	if err = filter.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	total, err := h.service.CalculateTotal(r.Context(), &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]int{"total": total}, http.StatusOK)
}

// Get godoc
// @Summary Информация о подписке
// @Description Посмотреть информацию о подписке по ее Id
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 200 {object} domain.Subscription
// @Failure 400 {string} string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	sub, err := h.service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sub == nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	writeJSON(w, sub, http.StatusOK)
}

// Update godoc
// @Summary Изменение записи подписки
// @Description Внесение изменений в подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки"
// @Param subscription body dto.UpdateSubscriptionRequest true "Тело запроса"
// @Success 200 {object} domain.Subscription
// @Failure 400 {string} string
// @Router /subscriptions/{id} [patch]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sub, err := h.service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sub == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if req.StartDate != nil {
		start, err := parseMonthYear(*req.StartDate)
		if err != nil {
			http.Error(w, "invalid start_date format", http.StatusBadRequest)
			return
		}
		sub.StartDate = start
	}

	if req.EndDate != nil {
		end, err := parseMonthYearPtr(req.EndDate)
		if err != nil {
			http.Error(w, "invalid end_date format", http.StatusBadRequest)
			return
		}

		sub.EndDate = end
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}

	if req.Price != nil {
		sub.Price = *req.Price
	}

	if err := h.service.Update(r.Context(), sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sub, http.StatusOK)
}

// Delete godoc
// @Summary Удаление записи о подписке
// @Description Удалить запись о подписке по её ID
// @Tags subscriptions
// @Param id path string true "UUID подписки"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary Список подписок
// @Description Получить список подписок (можно фильтровать по userId, serviceName)
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Наименование сервиса в подписке"
// @Success 200 {object} []domain.Subscription
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/list [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	userUuid := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	var userID *uuid.UUID
	if userUuid != "" {
		uid, err := parseUUID(userUuid)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		userID = &uid
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	list, err := h.service.List(r.Context(), userID, serviceNamePtr)
	if err != nil {
		http.Error(w, "Get subscriptions list error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, list, http.StatusOK)
}
