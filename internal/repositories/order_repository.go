package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodeyeen/wb-l0/internal/models"
	"github.com/kodeyeen/wb-l0/pkg/caching"
)

type OrderRepository interface {
	Add(context.Context, *models.Order) error
	GetByUid(context.Context, string) (*models.Order, error)
	GetAll(context.Context) ([]models.Order, error)
}

type orderRepository struct {
	dbpool *pgxpool.Pool
	cache  *caching.Cache
}

func NewOrderRepository(dbpool *pgxpool.Pool, cache *caching.Cache) *orderRepository {
	return &orderRepository{
		dbpool: dbpool,
		cache:  cache,
	}
}

func (or *orderRepository) Add(ctx context.Context, order *models.Order) error {
	query := `INSERT INTO orders (order_uid, data) VALUES ($1, $2)`

	_, err := or.dbpool.Exec(ctx, query, order.OrderUid, order.Data)

	if err != nil {
		return err
	}

	or.cache.Set(order.OrderUid, []byte(order.Data), 0)

	return err
}

func (or *orderRepository) GetByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	order := models.Order{}
	orderData, found := or.cache.Get(orderUid)

	if found {
		// log.Println("retrieved order from cache")

		order.OrderUid = orderUid
		order.Data = orderData.([]byte)

		return &order, nil
	}

	query := `SELECT order_uid, data FROM orders WHERE order_uid = $1`
	row := or.dbpool.QueryRow(ctx, query, orderUid)

	err := row.Scan(&order.OrderUid, &order.Data)

	if err != nil {
		return &models.Order{}, err
	}

	return &order, nil
}

func (or *orderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	query := `SELECT order_uid, data FROM orders`
	rows, _ := or.dbpool.Query(ctx, query)

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])

	if err != nil {
		return []models.Order{}, err
	}

	return orders, nil
}

func (or *orderRepository) RestoreCache(ctx context.Context) error {
	orders, err := or.GetAll(ctx)

	if err != nil {
		return err
	}

	for _, order := range orders {
		// go or.cache.Set(order.OrderUid, []byte(order.Data), 0)
		or.cache.Set(order.OrderUid, []byte(order.Data), 0)
	}

	log.Println("ORDERS", or.cache)

	return nil
}
