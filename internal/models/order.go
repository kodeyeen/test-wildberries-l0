package models

import "encoding/json"

type Order struct {
	OrderUid string          `db:"order_uid"`
	Data     json.RawMessage `db:"data"`
}
