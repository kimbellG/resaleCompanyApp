package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(lib_db *sql.DB) *OrderRepository {
	err := dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS ProductOrder (
	 		id	SERIAL PRIMARY KEY UNIQUE,
			clientId INTEGER REFERENCES Client (id),
			manager INTEGER REFERENCES userInformation (id),
			date timestamp,
			quantity INT NOT NULL,
			status VARCHAR(100) NOT NULL
		);`)

	if err != nil {
		panic(fmt.Errorf("order table: %v", err))
	}

	err = dbutil.Create(lib_db,
		`CREATE TABLE IF NOT EXISTS Purchases (
		id SERIAL PRIMARY KEY,
		orderID INT REFERENCES ProductOrder (id),
		offerID INT REFERENCES Offer (id)
	);`)

	if err != nil {
		panic(fmt.Errorf("purchases table: %v", err))
	}

	return &OrderRepository{
		db: lib_db,
	}
}

func (o *OrderRepository) Add(ctx context.Context, order *models.Order) error {
	orderId, err := o.addOrderInformation(order)
	if err != nil {
		return fmt.Errorf("add order information: %v", err)
	}

	if err := o.addOfferInformation(orderId, order.Offers); err != nil {
		return fmt.Errorf("add offer informarion: %v", err)
	}

	return nil
}

func (o *OrderRepository) addOrderInformation(order *models.Order) (int, error) {
	stmt, err := o.db.Prepare("INSERT INTO ProductOrder (clientId, manager, date, quantity, status) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return -1, fmt.Errorf("prepare stmt: %v", err)
	}

	managerId, err := o.getManagerID(order.ManagerLogin)
	if err != nil {
		return -1, fmt.Errorf("get manager id: %v", err)
	}

	var id int
	if err := stmt.QueryRow(order.ClientId, managerId, order.OrderDate, order.Quantity, order.Status).Scan(&id); err != nil {
		return -1, fmt.Errorf("exec stmt: %v", err)
	}

	return id, nil
}

func (o *OrderRepository) getManagerID(login string) (int, error) {
	stmt, err := o.db.Prepare("SELECT id FROM userInformation WHERE login = $1")
	if err != nil {
		return -1, fmt.Errorf("prepare stmt: %v", err)
	}

	var id int
	if err := stmt.QueryRow(login).Scan(&id); err != nil {
		return -1, fmt.Errorf("query stmt: %v", err)
	}

	return id, nil
}

func (o *OrderRepository) addOfferInformation(order int, offerIDs []int) error {
	stmt, err := o.db.Prepare("INSERT INTO Purchases (orderID, offerID) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	for _, id := range offerIDs {
		if _, err := stmt.Exec(order, id); err != nil {
			return fmt.Errorf("exec stmt: %v", err)
		}
	}

	return nil
}

func (o *OrderRepository) Gets(ctx context.Context) ([]models.Order, error) {
	result, err := o.getOrderInformation("SELECT * FROM ProductOrder")
	if err != nil {
		return nil, fmt.Errorf("get order information: %v", err)
	}

	for i, order := range result {
		result[i].Offers, err = o.getOfferInformation(order.Id)
		if err != nil {
			return nil, fmt.Errorf("get offer information: %v", err)
		}
	}

	return result, nil
}

func (o *OrderRepository) getOrderInformation(query string, arg ...interface{}) ([]models.Order, error) {
	stmt, err := o.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	rows, err := stmt.Query(arg...)
	if err != nil {
		return nil, fmt.Errorf("query stmt: %v", err)
	}
	defer rows.Close()

	result := make([]models.Order, 0)
	for rows.Next() {
		tmp := models.Order{}
		managerId := 0
		if err := rows.Scan(&tmp.Id, &tmp.ClientId, &managerId, &tmp.OrderDate, &tmp.Quantity, &tmp.Status); err != nil {
			return nil, fmt.Errorf("scaning result: %v", err)
		}

		tmp.ManagerLogin, err = o.getLoginManagerById(managerId)
		if err != nil {
			return nil, fmt.Errorf("get login: %v", err)
		}

		result = append(result, tmp)
	}

	return result, nil
}

func (o *OrderRepository) getLoginManagerById(id int) (string, error) {
	stmt, err := o.db.Prepare("SELECT login FROM userInformation WHERE id = $1")
	if err != nil {
		return "", fmt.Errorf("prepare stmt: %v", err)
	}

	var login string
	if err := stmt.QueryRow(id).Scan(&login); err != nil {
		return "", fmt.Errorf("scan login %v: %v", id, err)
	}

	return login, nil
}

func (o *OrderRepository) getOfferInformation(offer int) ([]int, error) {
	stmt, err := o.db.Prepare("SELECT offerID FROM Purchases WHERE orderID = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	rows, err := stmt.Query(offer)
	if err != nil {
		return nil, fmt.Errorf("query stmt: %v", err)
	}

	result := make([]int, 0)
	for rows.Next() {
		var offerID int
		if err := rows.Scan(&offerID); err != nil {
			return nil, fmt.Errorf("scan element: %v", err)
		}

		result = append(result, offerID)
	}

	return result, nil
}

func (o *OrderRepository) GetInInterval(ctx context.Context, start, end string) ([]models.Order, error) {
	result, err := o.getOrderInformation("SELECT * FROM ProductOrder WHERE date BETWEEN $1 AND $2", start, end)
	if err != nil {
		return nil, fmt.Errorf("get order information: %v", err)
	}

	for i, order := range result {
		result[i].Offers, err = o.getOfferInformation(order.Id)
		if err != nil {
			return nil, fmt.Errorf("get offer information: %v", err)
		}
	}

	return result, nil
}

func (o *OrderRepository) UpdateStatus(ctx context.Context, id int, newStatus string) error {
	stmt, err := o.db.Prepare("UPDATE ProductOrder SET status=$1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if _, err := stmt.Exec(newStatus, id); err != nil {
		return fmt.Errorf("exec stmt: %v", err)
	}

	return nil
}

func (o *OrderRepository) Filter(ctx context.Context, key string, value interface{}) ([]models.Order, error) {
	result, err := *new([]models.Order), error(nil)
	switch v := value.(type) {
	case string:
		result, err = o.getOrderInformation(fmt.Sprintf("SELECT * FROM ProductOrder WHERE %v LIKE $1", key), fmt.Sprintf("%%%v%%", v))
	default:
		result, err = o.getOrderInformation(fmt.Sprintf("SELECT * FROM ProductOrder WHERE %v = $1", key), v)
	}

	if err != nil {
		return nil, fmt.Errorf("repo order: %v", err)
	}

	for i, order := range result {
		result[i].Offers, err = o.getOfferInformation(order.Id)
		if err != nil {
			return nil, fmt.Errorf("repo offer: %v", err)
		}
	}

	return result, nil
}
