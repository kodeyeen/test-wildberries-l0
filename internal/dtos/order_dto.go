package dtos

import "time"

type OrderDTO struct {
	OrderUid          string      `json:"order_uid" validate:"required"`
	TrackNumber       string      `json:"track_number" validate:"required"`
	Entry             string      `json:"entry" validate:"required"`
	Delivery          DeliveryDTO `json:"delivery" validate:"required"`
	Payment           PaymentDTO  `json:"payment" validate:"required"`
	Items             []ItemDTO   `json:"items" validate:"required"`
	Locale            string      `json:"locale" validate:"required"`
	InternalSignature string      `json:"internal_signature"`
	CustomerId        string      `json:"customer_id" validate:"required"`
	DeliveryService   string      `json:"delivery_service" validate:"required"`
	Shardkey          string      `json:"shardkey" validate:"required"`
	SmId              int         `json:"sm_id" validate:"required"`
	DateCreated       time.Time   `json:"date_created" validate:"required"`
	OofShard          string      `json:"oof_shard" validate:"required"`
}
