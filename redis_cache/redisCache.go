package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Nirss/users/repository"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func NewCache(port string) *Cache {
	address := fmt.Sprintf("localhost:%v", port)
	return &Cache{client: redis.NewClient(&redis.Options{Addr: address})}
}

func (r *Cache) GetUsers(ctx context.Context) ([]repository.Users, error) {
	redisValue, err := r.client.Get(ctx, "users").Result()
	if err != nil {
		return nil, err
	}
	var result []repository.Users
	err = json.Unmarshal([]byte(redisValue), &result)
	if err != nil {
		log.Println("value unmarshal error: ", err)
		r.client.Del(ctx, "users")
		return nil, err
	}
	return result, nil
}

func (r *Cache) SetUsers(ctx context.Context, users []repository.Users) error {
	data, err := json.Marshal(users)
	if err != nil {
		log.Println("marshal error: ", err)
	}
	return r.client.Set(ctx, "users", data, time.Minute).Err()
}
