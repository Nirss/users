package main

import (
	"log"

	"github.com/Nirss/users/grpc"
	"github.com/Nirss/users/redis_cache"

	"github.com/Nirss/users/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	connect, err := ConnectDB()
	if err != nil {
		log.Println("connection database error: ", err)
	}
	var userRepository = repository.NewRepository(connect)
	var cache = redis_cache.NewCache("6379")
	grpc.ListenAndServe("localhost:8081", userRepository, cache)
}

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password='' dbname=teeest port=5432"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&repository.Users{})
	if err != nil {
		return nil, err
	}
	return db, err
}
