package grpc

import (
	"context"
	"log"
	"testing"

	"github.com/Nirss/users/redis_cache"

	grpcserver "github.com/Nirss/users/grpc/proto"
	"github.com/Nirss/users/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func Test_AddUsers(t *testing.T) {
	tests := []struct {
		name    string
		request *grpcserver.AddUserRequest
		err     error
	}{
		{
			name: "success",
			request: &grpcserver.AddUserRequest{
				Name:  "test",
				Email: "test",
				Phone: "123",
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := "host=localhost user=postgres password='' dbname=users port=5432"
			db, err := gorm.Open(postgres.New(postgres.Config{
				DSN: dsn,
			}), &gorm.Config{})
			if err != nil {
				log.Println("connection database error: ", err)
			}
			var userRepository = repository.NewRepository(db)
			service := &Service{usersRepository: userRepository}
			_, err = service.CreateUser(context.Background(), tt.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_DeleteUsers(t *testing.T) {
	tests := []struct {
		name    string
		request *grpcserver.DeleteUserRequest
		err     error
	}{
		{
			name: "success",
			request: &grpcserver.DeleteUserRequest{
				Id: 1,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := "host=localhost user=postgres password='' dbname=users port=5432"
			db, err := gorm.Open(postgres.New(postgres.Config{
				DSN: dsn,
			}), &gorm.Config{})
			if err != nil {
				log.Println("connection database error: ", err)
			}
			var userRepository = repository.NewRepository(db)
			service := &Service{usersRepository: userRepository}
			_, err = service.DeleteUser(context.Background(), tt.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_GetUsers(t *testing.T) {
	tests := []struct {
		name     string
		request  *grpcserver.GetUsersRequest
		wantBody []*grpcserver.User
		err      error
	}{
		{
			name:     "success",
			request:  &grpcserver.GetUsersRequest{},
			wantBody: []*grpcserver.User{},
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cache = redis_cache.NewCache("6379")
			dsn := "host=localhost user=postgres password='' dbname=users port=5432"
			db, err := gorm.Open(postgres.New(postgres.Config{
				DSN: dsn,
			}), &gorm.Config{})
			if err != nil {
				log.Println("connection database error: ", err)
			}
			var userRepository = repository.NewRepository(db)
			service := &Service{cache: cache, usersRepository: userRepository}
			users, err := service.GetAllUsers(context.Background())
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.wantBody, users)
		})
	}
}
