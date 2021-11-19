package store

import (
	"context"
	"6/internal/models"
)

type Store interface {

	Connect(url string) error
	Close() error

	
	Clothings() ClothingsRepository
	Accessories() AccessoriesRepository
	Users() UsersRepository
}


type ClothingsRepository interface {
	Create(ctx context.Context, clothing *models.Clothing) error
	All(ctx context.Context, filter *models.ClothingsFilter) ([]*models.Clothing, error)
	ByID(ctx context.Context, id int) (*models.Clothing, error)
	Update(ctx context.Context, clothing *models.Clothing) error
	Delete(ctx context.Context, id int) error
}

type AccessoriesRepository interface {
	Create(ctx context.Context, accessory *models.Accessory) error
	All(ctx context.Context, filter *models.AccessoriesFilter) ([]*models.Accessory, error)
	ByID(ctx context.Context, id int) (*models.Accessory, error)
	Update(ctx context.Context, accessory *models.Accessory) error
	Delete(ctx context.Context, id int) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context, filter *models.UsersFilter) ([]*models.User, error)
	ByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}