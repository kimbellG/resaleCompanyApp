package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type OfferRepository struct {
	db *sql.DB
}

func NewOfferRepository(lib_db *sql.DB) *OfferRepository {
	err := dbutil.CreateTable(lib_db,
		`CREATE TABLE IF NOT EXISTS Offer (
			id SERIAL PRIMARY KEY UNIQUE,
			productId INT REFERENCES Product (id),
			providerId INT REFERENCES Provider (vendor_code),
			cost DOUBLE PRECISION
		);`)

	if err != nil {
		panic(err)
	}

	return &OfferRepository{
		db: lib_db,
	}
}

func (r *OfferRepository) Add(ctx context.Context, offer *models.Offer) error {
	stmt, err := r.db.Prepare("INSERT INTO Offer (productId, providerId, cost) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("prepare statment: %v", err)
	}

	if _, err := stmt.Exec(offer.ProductId, offer.ProviderId, offer.Cost); err != nil {
		return fmt.Errorf("exec statment: %v", err)
	}

	return nil
}

func (r *OfferRepository) GetOfferOfProvider(ctx context.Context, providerId int) ([]models.Offer, error) {
	stmt, err := r.db.Prepare("SELECT * FROM Offer WHERE providerId = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	query, err := stmt.Query(providerId)
	if err != nil {
		return nil, fmt.Errorf("query stmt: %v", err)
	}
	defer query.Close()

	return scanOffer(query)
}

func scanOffer(rows *sql.Rows) ([]models.Offer, error) {
	offer := make([]models.Offer, 0)
	for rows.Next() {
		tmp := models.Offer{}
		if err := rows.Scan(&tmp.Id, &tmp.ProductId, &tmp.ProviderId, &tmp.Cost); err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		offer = append(offer, tmp)
	}

	return offer, nil
}

func (o *OfferRepository) GetOfferForProduct(ctx context.Context, productId int) ([]models.Offer, error) {
	stmt, err := o.db.Prepare("SELECT * FROM Offer WHERE productId = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	query, err := stmt.Query(productId)
	if err != nil {
		return nil, fmt.Errorf("exec stmt: %v", err)
	}
	defer query.Close()

	return scanOffer(query)
}

func (o *OfferRepository) UpdateCost(ctx context.Context, providerId, productId int, cost float32) error {
	stmt, err := o.db.Prepare("UPDATE Offer SET cost=$1 WHERE providerId = $2 AND productId = $3")
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if _, err := stmt.Exec(cost, providerId, productId); err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

func (o *OfferRepository) Delete(ctx context.Context, providerId, productId int) error {
	stmt, err := o.db.Prepare("DELETE FROM Offer WHERE providerId = $1 AND productId = $2")
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if _, err := stmt.Exec(providerId, productId); err != nil {
		return fmt.Errorf("exec stmt: %v", err)
	}

	return nil
}
