package database

import (
	"L0/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Repository interface {
	Save(order model.Order) (id int, err error)
	LoadToCache() (cache map[int]model.Order, err error)
}

type Config struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	DBName   string `mapstructure:"DB_NAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	SSLMode  string `mapstructure:"SSL_MODE"`
	RestPort string `mapstructure:"REST_PORT"`
}

type DataBase struct {
	client *sqlx.DB
}

func NewDataBase(cfg Config) (*DataBase, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode,
	)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return &DataBase{}, fmt.Errorf("could not connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
	}

	return &DataBase{client: db}, nil
}

func (d *DataBase) Save(order model.Order) (int, error) {
	d.client.MustExec("BEGIN;")
	defer d.client.Exec("ROLLBACK")
	d.client.Begin()

	query := "INSERT INTO delivery_desk (name, phone, zip, city, address, region, email) VALUES (:name, :phone, :zip, :city, :address, :region, :email) RETURNING id;"

	rows, err := d.client.NamedQuery(query, order.Delivery)
	if err != nil {
		return 0, fmt.Errorf("cannot insert delivery data: %w", err)
	}
	defer rows.Close()

	var deliveryID int

	for rows.Next() {
		if err := rows.Scan(&deliveryID); err != nil {
			return 0, fmt.Errorf("cannot scan deliveryID from query: %w", err)
		}
	}

	query = "INSERT INTO payment_desk (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES (:transaction, :request_id, :currency, :provider, :amount,:payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee) RETURNING id;"

	rows, err = d.client.NamedQuery(query, order.Payment)
	if err != nil {
		return 0, fmt.Errorf("cannot insert payment data: %w", err)
	}
	defer rows.Close()

	var paymentID int

	for rows.Next() {
		if err = rows.Scan(&paymentID); err != nil {
			return 0, fmt.Errorf("cannot scan PaymentID from query: %w", err)
		}
	}

	order.Payment.ID = paymentID
	order.Delivery.ID = deliveryID

	query = "INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id,locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)VALUES (:order_uid, :track_number, :entry, :d.id, :p.id,:locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard) RETURNING id;"

	rows, err = d.client.NamedQuery(query, order)
	if err != nil {
		return 0, fmt.Errorf("cannot insert order data: %w", err)
	}

	var orderID int

	for rows.Next() {
		if err = rows.Scan(&orderID); err != nil {
			return 0, fmt.Errorf("cannot scan orderID: %w", err)
		}
	}

	for _, item := range order.Items {
		item.OrderId = orderID

		query = "INSERT INTO items (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES (:order_id, :chrt_id, :track_number, :price, :rid, :name, :sale,:size, :total_price, :nm_id, :brand, :status)"

		_, err = d.client.NamedExec(query, item)
		if err != nil {
			return 0, fmt.Errorf("cannot insert item data: %w", err)
		}
	}

	d.client.MustExec("COMMIT")

	return orderID, nil
}

func (d *DataBase) LoadToCache() (cache map[int]model.Order, err error) {

	var amount int
	if err = d.client.QueryRow("SELECT COUNT(*) from orders").Scan(&amount); err != nil {
		return nil, fmt.Errorf("cannot get amount of orders from database: %w", err)
	}
	cache = make(map[int]model.Order, amount)

	query := "SELECT orders.id, order_uid, track_number, entry, locale, internal_signature, customer_id, " +
		"delivery_service, shardkey, sm_id, date_created, oof_shard, d.name as \"d.name\", d.phone as \"d.phone\", " +
		"d.zip as \"d.zip\", d.city as \"d.city\", d.address as \"d.address\", d.region as \"d.region\", " +
		"d.email as \"d.email\", p.transaction as \"p.transaction\", p.request_id as \"p.request_id\", " +
		"p.currency as \"p.currency\", p.provider as \"p.provider\", p.amount as \"p.amount\", " +
		"p.payment_dt as \"p.payment_dt\", p.bank as \"p.bank\", p.delivery_cost as \"p.delivery_cost\", " +
		"p.goods_total as \"p.goods_total\", p.custom_fee as \"p.custom_fee\" FROM orders " +
		"JOIN delivery_desk as d ON orders.delivery_id=d.id JOIN payment_desk as p ON orders.payment_id=p.id"

	rows, err := d.client.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("cannot get orders data: %w", err)
	}
	defer rows.Close()

	var nextOrder model.Order

	for rows.Next() {
		if err = rows.StructScan(&nextOrder); err != nil {
			return nil, fmt.Errorf("cannot parse order: %w", err)
		}
		cache[nextOrder.ID] = nextOrder
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in scanning order: %w", err)
	}

	query = "SELECT order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items"

	rows, err = d.client.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("cannot get items data: %w", err)
	}

	var nextItem model.Items

	for rows.Next() {
		if err = rows.StructScan(&nextItem); err != nil {
			return nil, fmt.Errorf("cannot parse item: %w", err)
		}

		var tmp = cache[nextItem.OrderId]
		tmp.Items = append(cache[nextItem.OrderId].Items, nextItem)
		cache[nextItem.OrderId] = tmp

	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in scanning items: %w", err)
	}

	return cache, nil

}
