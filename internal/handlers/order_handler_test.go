package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kodeyeen/wb-l0/internal/dtos"
	"github.com/kodeyeen/wb-l0/internal/mocks"
	"github.com/kodeyeen/wb-l0/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderHandler(t *testing.T) {
	uid, _ := uuid.NewRandom()
	orderUid := uid.String()

	expectedResp := []byte(fmt.Sprintf(`{
		"order_uid": "%s",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
			"name": "Test Testov",
			"phone": "+9720000000",
			"zip": "2639809",
			"city": "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region": "Kraiot",
			"email": "test@gmail.com"
		},
		"payment": {
			"transaction": "%s",
			"request_id": "",
			"currency": "USD",
			"provider": "wbpay",
			"amount": 1817,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 317,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 9934930,
				"track_number": "WBILMTESTTRACK",
				"price": 453,
				"rid": "ab4219087a764ae0btest",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 202
			}
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	}`, orderUid, orderUid))

	mockOrder := &models.Order{
		OrderUid: orderUid,
		Data:     expectedResp,
	}

	mockOrderService := new(mocks.MockOrderService)
	orderHandler := NewOrderHandler(mockOrderService)

	mockOrderService.On("GetByUid", mock.Anything, orderUid).Return(mockOrder, nil)

	url := fmt.Sprintf("/order?order_uid=%s", orderUid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	orderHandler.GetOrderByUid(recorder, req)

	var orderDTO dtos.OrderDTO
	err = json.Unmarshal(recorder.Body.Bytes(), &orderDTO)

	mockOrderService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status 200 but got %d", recorder.Code)
	assert.NoError(t, err, "error when unmarshalling response body")
}
