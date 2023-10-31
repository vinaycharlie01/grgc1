package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	userpb "myapp/user/proto"
)

type UserServiceServer struct {
	users                                 map[int32]*userpb.User
	userpb.UnimplementedUserServiceServer // Embed the UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetUserById(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, exists := s.users[req.UserId]
	if !exists {
		return nil, fmt.Errorf("User with ID %d not found", req.UserId)
	}
	return user, nil
}

func (s *UserServiceServer) GetUsersByIds(req *userpb.GetUsersRequest, stream userpb.UserService_GetUsersByIdsServer) error {
	for _, id := range req.UserIds {
		user, exists := s.users[id]
		if exists {
			if err := stream.Send(user); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userService := &UserServiceServer{
		users: map[int32]*userpb.User{
			1: {
				Id:      1,
				Fname:   "Steve",
				City:    "LA",
				Phone:   1234567890,
				Height:  5.8,
				Married: true,
			},
			2: {
				Id:      2,
				Fname:   "Steve",
				City:    "LA",
				Phone:   1234567890,
				Height:  5.8,
				Married: true,
			},
			// Add more users as needed
		},
	}
	userpb.RegisterUserServiceServer(s, userService)
	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
