package redis_cache

import (
	"context"
	"encoding/json"
	"6/internal/cache"
	"6/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

func (rc RedisCache) Accessories() cache.AccessoriesCacheRepo {
	if rc.accessories == nil {
		rc.accessories = newAccessoriesRepo(rc.client, rc.expires)
	}

	return rc.accessories
}


type AccessoriesRepo struct {
	client  *redis.Client
	expires time.Duration
}
}

func newAccessoriesRepo(client *redis.Client, exp time.Duration) cache.AccessoriesCacheRepo {
	return &AccessoriesRepo{
		client:  client,
		expires: exp,
	}
}



func (c AccessoriesRepo) Set(ctx context.Context, key string, value []*models.Accessory) error {
	accessoryBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, accessoryBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}


func (c AccessoriesRepo) Get(ctx context.Context, key string) ([]*models.Accessory, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	accessories := make([]*models.Accessory, 0)
	if err = json.Unmarshal([]byte(result), &accessories); err != nil {
		return nil, err
	}

	return accessories, nil
}
