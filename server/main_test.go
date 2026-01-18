package main

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/steviepreston/grpc-proto-example/proto-buffers/users/v1"
)

func TestGetUser_EmptyID(t *testing.T) {
	s := newUserServer()
	_, err := s.GetUser(context.Background(), &pb.GetUserRequest{Id: ""})
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestGetUser_NotFound(t *testing.T) {
	s := newUserServer()
	_, err := s.GetUser(context.Background(), &pb.GetUserRequest{Id: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("expected NotFound, got %v", status.Code(err))
	}
}

func TestGetUser_Success(t *testing.T) {
	s := newUserServer()
	createResp, _ := s.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	})

	resp, err := s.GetUser(context.Background(), &pb.GetUserRequest{Id: createResp.User.Id})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.User.Email)
	}
}

func TestCreateUser_MissingEmail(t *testing.T) {
	s := newUserServer()
	_, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email: "",
		Name:  "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing email")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	s := newUserServer()
	_, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "",
	})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestCreateUser_Success(t *testing.T) {
	s := newUserServer()
	resp, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email: "alice@example.com",
		Name:  "Alice",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Id == "" {
		t.Error("expected user ID to be set")
	}
	if resp.User.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", resp.User.Email)
	}
	if resp.User.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", resp.User.Name)
	}
	if resp.User.CreatedAt == 0 {
		t.Error("expected CreatedAt to be set")
	}
}

func TestListUsers_Empty(t *testing.T) {
	s := newUserServer()
	resp, err := s.ListUsers(context.Background(), &pb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Users) != 0 {
		t.Errorf("expected 0 users, got %d", len(resp.Users))
	}
}

func TestListUsers_WithUsers(t *testing.T) {
	s := newUserServer()
	s.CreateUser(context.Background(), &pb.CreateUserRequest{Email: "a@example.com", Name: "A"})
	s.CreateUser(context.Background(), &pb.CreateUserRequest{Email: "b@example.com", Name: "B"})
	s.CreateUser(context.Background(), &pb.CreateUserRequest{Email: "c@example.com", Name: "C"})

	resp, err := s.ListUsers(context.Background(), &pb.ListUsersRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Users) != 3 {
		t.Errorf("expected 3 users, got %d", len(resp.Users))
	}
}

type mockStreamUsersServer struct {
	pb.UserService_StreamUsersServer
	users []*pb.StreamUsersResponse
}

func (m *mockStreamUsersServer) Send(resp *pb.StreamUsersResponse) error {
	m.users = append(m.users, resp)
	return nil
}

func TestStreamUsers_Empty(t *testing.T) {
	s := newUserServer()
	mock := &mockStreamUsersServer{}

	err := s.StreamUsers(&pb.StreamUsersRequest{}, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.users) != 0 {
		t.Errorf("expected 0 streamed users, got %d", len(mock.users))
	}
}

func TestStreamUsers_WithUsers(t *testing.T) {
	s := newUserServer()
	s.CreateUser(context.Background(), &pb.CreateUserRequest{Email: "a@example.com", Name: "A"})
	s.CreateUser(context.Background(), &pb.CreateUserRequest{Email: "b@example.com", Name: "B"})

	mock := &mockStreamUsersServer{}
	err := s.StreamUsers(&pb.StreamUsersRequest{}, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.users) != 2 {
		t.Errorf("expected 2 streamed users, got %d", len(mock.users))
	}
}
