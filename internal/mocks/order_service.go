package mocks

import (
	"context"

	"github.com/kodeyeen/wb-l0/internal/models"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) GetByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	args := m.Called(ctx, orderUid)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) HandleMessage(msg *stan.Msg) {}

func (m *MockOrderService) Save(data []byte) (*models.Order, error) {
	args := m.Called(data)
	return args.Get(0).(*models.Order), args.Error(1)
}
