package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userpb "myapp/user/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	srv := grpc.NewServer()

	client := userpb.NewUserServiceClient(conn)
	reflection.Register(srv)
	// Fetch a single user by ID

	mux := http.NewServeMux()
	mux.HandleFunc("/GetUserById", func(w http.ResponseWriter, r *http.Request) {
		num1, _ := strconv.Atoi(r.FormValue("Id"))

		responce, err := client.GetUserById(context.Background(), &userpb.GetUserRequest{UserId: int32(num1)})
		if err != nil {
			log.Fatalf("Error while calling GetUserById: %v", err)
		}
		fmt.Fprintf(w, "Result: %v\n", responce)
	})
	fmt.Println("Server Started at port 3333")
	http.ListenAndServe(":3333", mux)

	// fmt.Println("User Details (GetUserById):", getUserResponse)

	// Fetch multiple users by a list of IDs
	getUsersResponse, err := client.GetUsersByIds(context.Background(), &userpb.GetUsersRequest{UserIds: []int32{1, 2, 3}})
	if err != nil {
		log.Fatalf("Error while calling GetUsersByIds: %v", err)
	}
	for {
		_, err := getUsersResponse.Recv()
		if err != nil {
			break
		}
		// continue
		// fmt.Println("User Details (GetUsersByIds):", user)
	}
}
