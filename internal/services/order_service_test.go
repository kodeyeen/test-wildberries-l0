package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/kodeyeen/wb-l0/internal/mocks"
	"github.com/kodeyeen/wb-l0/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func genOrderData(orderUid string) []byte {
	return []byte(fmt.Sprintf(`{
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
}

func TestOrderService_GetByUid(t *testing.T) {
	uid, _ := uuid.NewRandom()
	orderUid := uid.String()
	ctx := context.Background()

	mockRepository := new(mocks.MockOrderRepository)
	orderService := NewOrderService(mockRepository)

	expectedOrder := &models.Order{
		OrderUid: orderUid,
		Data:     genOrderData(orderUid),
	}
	mockRepository.On("GetByUid", ctx, orderUid).Return(expectedOrder, nil)

	retrievedOrder, err := orderService.GetByUid(ctx, orderUid)

	mockRepository.AssertExpectations(t)
	assert.NoError(t, err, "error when retrieving order")
	assert.Equal(t, expectedOrder, retrievedOrder, "orders don't match")
}

func TestOrderService_Save(t *testing.T) {
	uid, _ := uuid.NewRandom()
	orderUid := uid.String()
	orderData := genOrderData(orderUid)

	mockRepository := new(mocks.MockOrderRepository)
	orderService := NewOrderService(mockRepository)

	expectedOrder := &models.Order{
		OrderUid: orderUid,
		Data:     orderData,
	}
	mockRepository.On("Add", mock.Anything, expectedOrder).Return(expectedOrder, nil)
	savedOrder, err := orderService.Save(orderData)

	mockRepository.AssertExpectations(t)
	assert.NoError(t, err, "error when saving order")
	assert.Equal(t, expectedOrder, savedOrder, "orders don't match")
}
