package mocks

import (
	"context"

	"github.com/kodeyeen/wb-l0/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Add(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(1)
}

func (m *MockOrderRepository) GetByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	args := m.Called(ctx, orderUid)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Order), args.Error(1)
}
