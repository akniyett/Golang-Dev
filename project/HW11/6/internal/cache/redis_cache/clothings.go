package redis_cache

import (
	"context"
	"encoding/json"
	"6/internal/cache"
	"6/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)
func (rc RedisCache) Clothings() cache.ClothingsCacheRepo {
	if rc.clothings == nil {
		rc.clothings = newClothingsRepo(rc.client, rc.expires)
	}

	return rc.clothings
}
type ClothingsRepo struct {
	client  *redis.Client
	expires time.Duration
}
func newClothingsRepo(client *redis.Client, exp time.Duration) cache.ClothingsCacheRepo {
	return &ClothingsRepo{
		client:  client,
		expires: exp,
	}
}
func (c ClothingsRepo) Set(ctx context.Context, key string, value []*models.Clothing) error {
	clothingBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, clothingBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}
func (c ClothingsRepo) Get(ctx context.Context, key string) ([]*models.Clothing, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	clothings := make([]*models.Accessory, 0)
	if err = json.Unmarshal([]byte(result), &clothings); err != nil {
		return nil, err
	}

	return clothings, nil
}