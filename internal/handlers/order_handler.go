package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kodeyeen/wb-l0/internal/services"
)

type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

func (oh *orderHandler) GetOrderByUid(w http.ResponseWriter, r *http.Request) {
	order_uid := r.URL.Query().Get("order_uid")
	order, err := oh.orderService.GetByUid(r.Context(), order_uid)

	if err != nil {
		log.Println("err", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order.Data)
}
