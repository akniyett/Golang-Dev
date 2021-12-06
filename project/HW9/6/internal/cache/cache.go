package cache

import (
	"context"
	"6/internal/models"
)

type Cache interface {
	Close() error

	Accessories() AcessoriesCacheRepo
	Users() UsersCacheRepo
	Clothings() ClothingsCacheRepo

	DeleteAll(ctx context.Context) error
}

type AccessoriesCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Accessory) error
	Get(ctx context.Context, key string) ([]*models.Accessory, error)
}

type UsersCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.User) error
	Get(ctx context.Context, key string) ([]*models.User, error)
}

type ClothingsCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Clothing) error
	Get(ctx context.Context, key string) ([]*models.Clothing, error)
}