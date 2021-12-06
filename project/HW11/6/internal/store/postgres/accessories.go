package postgres

import (
	"context"
	"6/internal/models"
	"6/internal/store"
	"github.com/jmoiron/sqlx"
	
)



func (db *DB) Accessories() store.AccessoriesRepository {
	if db.accessories == nil {
		db.accessories = NewAccessoriesRepository(db.conn)
	}

	return db.accessories
}

type AccessoriesRepository struct {
	conn *sqlx.DB
}

func NewAccessoriesRepository(conn *sqlx.DB) store.AccessoriesRepository {
	return &AccessoriesRepository{conn: conn}
}

func (a *AccessoriesRepository) Create(ctx context.Context, accessory *models.Accessory) error {
	_, err := a.conn.Exec("INSERT INTO accessories(name, description, size, price, isAvailable, manufacturer, material) VALUES ($1, $2, %3, $4, $5, $6, $7)", accessory.Name, accessory.Description, accessory.Size, accessory.Price, accessory.IsAvailable, accessory.Manufacturer, accessory.Material)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccessoriesRepository) All(ctx context.Context, filter *models.AccessoriesFilter) ([]*models.Accessory, error) {
	

	accessories := make([]*models.Accessory, 0)
	basicQuery := "SELECT * FROM accessories"
	if filter.Query != nil {
		// basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)
		basicQuery += " WHERE name ILIKE '%" + *filter.Query + "%'"
		// if err := a.conn.Select(&accessories, basicQuery, "%"+*filter.Query+"%"); err != nil {
		// 	return nil, err
		// }
	}
	if err := a.conn.SelectContext(ctx, &accessories,  basicQuery); err != nil {
		return nil, err
	}

	return accessories, nil
}

func (a *AccessoriesRepository) ByID(ctx context.Context, id int) (*models.Accessory, error) {

	accessory := new(models.Accessory)
	if err := a.conn.GetContext(ctx, accessory, "SELECT id, name, description, size, price, isAvailable, manufacturer, material FROM accessories WHERE id=$1", id); err != nil {
		return nil, err
	}

	return accessory, nil
}

func (a *AccessoriesRepository) Update(ctx context.Context, accessory *models.Accessory) error {
	_, err := a.conn.ExecContext(ctx, "UPDATE accessories SET name = $1, description = $2, size = $3, price = $4, isAvailable = $5, manufacturer = $6, material = $7 WHERE id = $8", accessory.Name, accessory.Description, accessory.Size, accessory.Price, accessory.IsAvailable, accessory.Manufacturer, accessory.Material , accessory.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a AccessoriesRepository) Delete(ctx context.Context, id int) error {
	_, err := a.conn.ExecContext("DELETE FROM accessories WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
