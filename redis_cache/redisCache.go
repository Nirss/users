package redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Nirss/users/repository"

	"github.com/go-redis/redis/v8"
)

var (
	ErrUnexpectedError = errors.New("unexpected error, please repeat again later")
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
	var result []repository.Users
	if err != nil {
		log.Println("get redis value error: ", err)
		return nil, ErrUnexpectedError
	}
	err = json.Unmarshal([]byte(redisValue), &result)
	if err != nil {
		log.Println("get redis value error: ", err)
		r.client.Del(ctx, "users")
		return nil, ErrUnexpectedError
	}
	return result, nil
}

func (r *Cache) SetUsers(ctx context.Context, users []repository.Users) error {
	data, err := json.Marshal(users)
	if err != nil {
		log.Println("set redis value error: ", err)
		return ErrUnexpectedError
	}
	err = r.client.Set(ctx, "users", data, 0).Err()
	if err != nil {
		log.Println("set redis value error: ", err)
	}
	return nil
}
