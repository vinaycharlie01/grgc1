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

// / Create a gRPC method to receive and process user data
func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.User) (*userpb.User, error) {
	// Process the user data and store it as needed
	// You can create a unique ID or use your own logic

	// For example, you can generate a unique ID:
	// userID := GenerateUniqueUserID()

	// Create a new user based on the request data
	newUser := &userpb.User{
		Id:      req.Id,
		Fname:   req.Fname,
		City:    req.City,
		Phone:   req.Phone,
		Height:  req.Height,
		Married: req.Married,
	}

	// Store the new user in your data structure (e.g., a map)
	s.users[req.Id] = newUser

	return newUser, nil
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
