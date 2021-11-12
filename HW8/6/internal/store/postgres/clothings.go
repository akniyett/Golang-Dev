package postgres

import (
	"context"

	"6/internal/models"
	"6/internal/store"
	"github.com/jmoiron/sqlx"
	
)

func (db *DB) Clothings() store.ClothingsRepository {
	if db.clothings == nil {
		db.clothings = NewClothingsRepository(db.conn)
	}

	return db.clothings
}

type ClothingsRepository struct {
	conn *sqlx.DB
}

func NewClothingsRepository(conn *sqlx.DB) store.ClothingsRepository {
	return &ClothingsRepository{conn: conn}
}

func (c ClothingsRepository) Create(ctx context.Context, clothing *models.Clothing) error {
	_, err := c.conn.Exec("INSERT INTO clothings(name, description, size, price, isAvailable) VALUES ($1, $2, %3, $4, $5)", clothing.Name, clothing.Description, clothing.Size, clothing.Price, clothing.IsAvailable)
	if err != nil {
		return err
	}
	return nil
}

func (c ClothingsRepository) All(ctx context.Context) ([]*models.Clothing, error) {
	clothings := make([]*models.Clothing, 0)
	if err := c.conn.Select(&clothings, "SELECT * FROM clothings"); err != nil {
		return nil, err
	}

	return clothings, nil
}

func (c ClothingsRepository) ByID(ctx context.Context, id int) (*models.Clothing, error) {
	clothing := new(models.Clothing)
	if err := c.conn.Get(clothing, "SELECT id, name, description, size, price, isAvailable FROM accessories WHERE id=$1", id); err != nil {
		return nil, err
	}


	return clothing, nil
}

func (c ClothingsRepository) Update(ctx context.Context, clothing *models.Clothing) error {
	_, err := c.conn.Exec("UPDATE clothings SET name = $1, description = $2, size = $3, price = $4, isAvailable = $5 WHERE id = $8", clothing.Name, clothing.Description, clothing.Size, clothing.Price, clothing.IsAvailable, clothing.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c ClothingsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM clothings WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}