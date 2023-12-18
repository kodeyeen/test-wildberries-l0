package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/kodeyeen/test-wildberries-l0/internal/dtos"
	"github.com/kodeyeen/test-wildberries-l0/internal/models"
	"github.com/kodeyeen/test-wildberries-l0/internal/repositories"
	"github.com/nats-io/stan.go"
)

type OrderService interface {
	GetByUid(context.Context, string) (*models.Order, error)
	HandleMessage(*stan.Msg)
	Save([]byte) (*models.Order, error)
}

type orderService struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) *orderService {
	return &orderService{
		orderRepository: orderRepository,
	}
}

func (os *orderService) GetByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	order, err := os.orderRepository.GetByUid(ctx, orderUid)

	if err != nil {
		return &models.Order{}, err
	}

	return order, nil
}

func (os *orderService) HandleMessage(msg *stan.Msg) {
	log.Printf("Received message: %s\n", string(msg.Data))

	os.Save(msg.Data)
}

func (os *orderService) Save(data []byte) (*models.Order, error) {
	var orderDTO dtos.OrderDTO
	err := json.Unmarshal(data, &orderDTO)

	if err != nil {
		log.Printf("error when unmarshaling message: %v", err)
		return &models.Order{}, nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(orderDTO)

	if err != nil {
		log.Printf("message validation failed: %v", err)
		return &models.Order{}, nil
	}

	order := &models.Order{
		OrderUid: orderDTO.OrderUid,
		Data:     data,
	}

	err = os.orderRepository.Add(context.Background(), order)

	if err != nil {
		log.Printf("error when saving order: %v", err)
	}

	return order, nil
}
