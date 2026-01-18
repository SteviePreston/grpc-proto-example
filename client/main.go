// client/main.go
package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/steviepreston/grpc-proto-example/proto-buffers/users/v1"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create users
	for _, name := range []string{"Alice", "Bob", "Charlie"} {
		resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
			Email: name + "@example.com",
			Name:  name,
		})
		if err != nil {
			log.Fatalf("CreateUser failed: %v", err)
		}
		log.Printf("Created user: %+v", resp.User)
	}

	// List users
	listResp, err := client.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	log.Printf("All users: %d", len(listResp.Users))

	// Stream users
	stream, err := client.StreamUsers(ctx, &pb.StreamUsersRequest{})
	if err != nil {
		log.Fatalf("StreamUsers failed: %v", err)
	}

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}
		log.Printf("Streamed user: %s", user.User.Name)
	}
}
