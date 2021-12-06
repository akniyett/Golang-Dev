package redis_cache

import (
	"context"
	"encoding/json"
	"6/internal/cache"
	"6/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

func (rc RedisCache) Users() cache.UsersCacheRepo {
	if rc.users == nil {
		rc.users = newUsersRepo(rc.client, rc.expires)
	}

	return rc.users
}


type UsersRepo struct {
	client  *redis.Client
	expires time.Duration
}
}

func newUsersRepo(client *redis.Client, exp time.Duration) cache.UsersCacheRepo {
	return &UsersRepo{
		client:  client,
		expires: exp,
	}
}



func (c UsersRepo) Set(ctx context.Context, key string, value []*models.User) error {
	userBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, userBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}


func (c UserRepo) Get(ctx context.Context, key string) ([]*models.User, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	users := make([]*models.User, 0)
	if err = json.Unmarshal([]byte(result), &users); err != nil {
		return nil, err
	}

	return users, nil
}
