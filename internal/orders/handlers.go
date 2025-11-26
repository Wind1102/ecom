package orders

import (
	"net/http"

	"github.com/wind1102/ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {

	}
	h.service.PlaceOrder(r.Context(), tempOrder)
	json.Write(w, http.StatusCreated, nil)
}
