package entity

import "time"

type Order struct {
	OrderUID          string    `db:"order_uid" json:"order_uid"`
	TrackNumber       string    `db:"track_number" json:"track_number"`
	Entry             string    `db:"entry" json:"entry"`
	Locale            string    `db:"locale" json:"locale"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature"`
	CustomerID        string    `db:"customer_id" json:"customer_id"`
	DeliveryService   string    `db:"delivery_service" json:"delivery_service"`
	ShardKey          string    `db:"shardkey" json:"shardkey"`
	SmID              int       `db:"sm_id" json:"sm_id"`
	DateCreated       time.Time `db:"date_created" json:"date_created"`
	OofShard          string    `db:"oof_shard" json:"oof_shard"`

	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Item   `json:"items"`
}

type Delivery struct {
	OrderUID string `db:"order_uid" json:"order_uid"`
	Name     string `db:"name" json:"name"`
	Phone    string `db:"phone" json:"phone"`
	Zip      string `db:"zip" json:"zip"`
	City     string `db:"city" json:"city"`
	Address  string `db:"address" json:"address"`
	Region   string `db:"region" json:"region"`
	Email    string `db:"email" json:"email"`
}

type Payment struct {
	Transaction  string  `db:"transaction" json:"transaction"`
	OrderUID     string  `db:"order_uid" json:"order_uid"`
	RequestID    string  `db:"request_id" json:"request_id"`
	Currency     string  `db:"currency" json:"currency"`
	Provider     string  `db:"provider" json:"provider"`
	Amount       float64 `db:"amount" json:"amount"`
	PaymentDt    int64   `db:"payment_dt" json:"payment_dt"`
	Bank         string  `db:"bank" json:"bank"`
	DeliveryCost float64 `db:"delivery_cost" json:"delivery_cost"`
	GoodsTotal   int     `db:"goods_total" json:"goods_total"`
	CustomFee    float64 `db:"custom_fee" json:"custom_fee"`
}

type Item struct {
	ChrtID      int     `db:"chrt_id" json:"chrt_id"`
	OrderUID    string  `db:"order_uid" json:"order_uid"`
	TrackNumber string  `db:"track_number" json:"track_number"`
	Price       float64 `db:"price" json:"price"`
	RID         string  `db:"rid" json:"rid"`
	Name        string  `db:"name" json:"name"`
	Sale        int     `db:"sale" json:"sale"`
	Size        string  `db:"size" json:"size"`
	TotalPrice  float64 `db:"total_price" json:"total_price"`
	NmID        int     `db:"nm_id" json:"nm_id"`
	Brand       string  `db:"brand" json:"brand"`
	Status      int     `db:"status" json:"status"`
}
