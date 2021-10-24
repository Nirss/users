package grpc

import (
	"context"
	"log"
	"net"

	"gorm.io/gorm"

	grpcserver "github.com/Nirss/users/grpc/proto"
	"github.com/Nirss/users/redis_cache"
	"github.com/Nirss/users/repository"
	"google.golang.org/grpc"
)

type Service struct {
	grpcserver.UnimplementedUsersServiceServer
	cache           *redis_cache.Cache
	usersRepository *repository.UsersRepository
}

func (s *Service) GetUsers(ctx context.Context, request *grpcserver.GetUsersRequest) (*grpcserver.GetUsersResponse, error) {
	result, err := s.cache.GetUsers(ctx)
	if err != nil {
		result, err = s.usersRepository.GetUsers()
		if err != nil {
			return nil, err
		}
		err = s.cache.SetUsers(context.Background(), result)
		if err != nil {
			log.Println("set redis value error: ", err)
		}
	}
	var users []*grpcserver.User
	for _, user := range result {
		users = append(users, &grpcserver.User{Name: user.Name, Phone: user.Phone, Email: user.Email, Id: user.Id})
	}
	return &grpcserver.GetUsersResponse{Result: users}, nil
}

func (s *Service) AddUser(ctx context.Context, request *grpcserver.AddUserRequest) (*grpcserver.AddUserResponse, error) {
	var params = &repository.Users{
		Model: gorm.Model{},
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}
	err := s.usersRepository.CreateUser(params)
	return &grpcserver.AddUserResponse{}, err
}

func (s *Service) DeleteUser(ctx context.Context, request *grpcserver.DeleteUserRequest) (*grpcserver.DeleteUserResponse, error) {
	err := s.usersRepository.DeleteUser(int(request.Id))
	return &grpcserver.DeleteUserResponse{}, err
}

func ListenAndServe(port string, usersRepository *repository.UsersRepository, cache *redis_cache.Cache) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	opt := []grpc.ServerOption{}
	server := grpc.NewServer(opt...)
	grpcserver.RegisterUsersServiceServer(server, &Service{usersRepository: usersRepository, cache: cache})
	server.Serve(listener)
}
