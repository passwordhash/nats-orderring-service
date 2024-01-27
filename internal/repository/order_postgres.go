package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"nats_server/internal/entity"
)

var OrderNotFoundErr = errors.New("order not found")
var AdditionsNotFoundErr = errors.New("some of addition inforamation about order not found")

type OrderRepository struct {
	psqlDb *sqlx.DB
}

func NewOrderPostgres(psqlDb *sqlx.DB) *OrderRepository {
	return &OrderRepository{psqlDb: psqlDb}
}

func (r *OrderRepository) Create(o entity.Order) (string, error) {
	tx, err := r.psqlDb.Beginx()
	if err != nil {
		return "", err
	}

	// new order
	_, err = tx.Exec(`
	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// new delivery
	_, err = tx.Exec(`
	INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip, o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec(`
	INSERT INTO payment (transaction, order_uid, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		o.Payment.Transaction, o.OrderUID, o.Payment.RequestID, o.Payment.Currency, o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	for _, item := range o.Items {
		_, err = tx.Exec(`
		INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			item.ChrtID, o.OrderUID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	return o.OrderUID, tx.Commit()
}

func (r *OrderRepository) GetWithAddition(orderUID string) (entity.Order, error) {
	var order entity.Order

	err := r.psqlDb.Get(&order, "SELECT * FROM orders WHERE order_uid = $1", orderUID)
	if err != nil {
		return order, OrderNotFoundErr
	}

	err = r.enrichOrder(&order)

	return order, err
}

func (r *OrderRepository) GetWithAdditionList() ([]entity.Order, error) {
	var list []entity.Order

	err := r.psqlDb.Select(&list, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	logrus.Info(list)

	for i := range list {
		err = r.enrichOrder(&list[i])
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

// enrichOrder enriches order with addition information
func (r *OrderRepository) enrichOrder(order *entity.Order) error {
	err := r.psqlDb.Get(&order.Delivery, "SELECT * FROM delivery WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return errors.Join(AdditionsNotFoundErr, err)
	}

	err = r.psqlDb.Get(&order.Payment, "SELECT * FROM payment WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return errors.Join(AdditionsNotFoundErr, err)
	}

	err = r.psqlDb.Select(&order.Items, "SELECT * FROM items WHERE order_uid = $1", order.OrderUID)

	return nil
}
