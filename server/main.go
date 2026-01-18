// server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/steviepreston/grpc-proto-example/proto-buffers"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	mu    sync.RWMutex
	users map[string]*pb.User
}

func newUserServer() *userServer {
	return &userServer{
		users: make(map[string]*pb.User),
	}
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{User: user}, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Email == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "email and name are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("user_%d", len(s.users)+1)
	user := &pb.User{
		Id:        id,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
	}
	s.users[id] = user

	return &pb.CreateUserResponse{User: user}, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*pb.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}

	return &pb.ListUsersResponse{Users: users}, nil
}

func (s *userServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if err := stream.Send(user); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond) // simulate delay
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, newUserServer())

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
