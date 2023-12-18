package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kodeyeen/wb-l0/internal/handlers"
	"github.com/kodeyeen/wb-l0/internal/repositories"
	"github.com/kodeyeen/wb-l0/internal/services"
	"github.com/kodeyeen/wb-l0/pkg/caching"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	loadDotEnv()

	ctx := context.Background()

	dbpool := connectToDB(ctx)
	natsConn := connectToNats()
	stanConn := connectToStan(natsConn)

	defer dbpool.Close()
	defer natsConn.Close()
	defer stanConn.Close()

	ordersCache := caching.New(10*time.Minute, 2*time.Minute)

	orderRepository := repositories.NewOrderRepository(dbpool, ordersCache)
	orderService := services.NewOrderService(orderRepository)
	orderHandler := handlers.NewOrderHandler(orderService)

	err := orderRepository.RestoreCache(ctx)
	if err != nil {
		log.Fatalf("failed to restore orders cache: %v", err)
	}

	_, err = stanConn.Subscribe("order:created", orderService.HandleMessage)
	if err != nil {
		log.Fatalf("unable to subscribe to channel: %v", err)
	}

	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.HandleFunc("/order", orderHandler.GetOrderByUid)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadDotEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("unable to load .env file: %v", err)
	}
}

func connectToDB(ctx context.Context) *pgxpool.Pool {
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to create connection pool: %v\n", err)
	}

	return dbpool
}

func connectToNats() *nats.Conn {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("unable to connect to NATS: %v\n", err)
	}

	return nc
}

func connectToStan(natsConn *nats.Conn) stan.Conn {
	sc, err := stan.Connect(
		os.Getenv("STAN_CLUSTER_ID"),
		os.Getenv("STAN_CLIENT_ID"),
		stan.NatsConn(natsConn),
	)
	if err != nil {
		log.Fatalf("unable to connect to NATS Streaming: %v\n", err)
	}

	return sc
}
