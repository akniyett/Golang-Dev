package inmemory

import (
	"context"
	"fmt"
	"6/internal/models"
	"6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Product

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Product),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, product *models.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[product.ID] = product
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	products := make([]*models.Product, 0, len(db.data))
	for _, product := range db.data {
		products = append(products, product)
	}

	return products, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	product, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No product with id %d", id)
	}

	return product, nil
}

func (db *DB) Update(ctx context.Context, product *models.Product) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[product.ID] = product
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
