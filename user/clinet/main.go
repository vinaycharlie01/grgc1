package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	userpb "myapp/user/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/GetUserById", func(w http.ResponseWriter, r *http.Request) {
		num1, _ := strconv.Atoi(r.FormValue("Id"))

		response, err := client.GetUserById(context.Background(), &userpb.GetUserRequest{UserId: int32(num1)})
		if err != nil {
			log.Fatalf("Error while calling GetUserById: %v", err)
		}
		fmt.Fprintf(w, "Result: %v\n", response)
	})

	mux.HandleFunc("/GetUsersByIds", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.FormValue("Ids")
		ids := strings.Split(idStr, ",")
		var userIDs []int32

		for _, id := range ids {
			num, err := strconv.Atoi(id)
			if err != nil {
				log.Fatalf("Error parsing ID: %v", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			userIDs = append(userIDs, int32(num))
		}

		// Call the gRPC server with the user IDs
		getUsersResponse, err := client.GetUsersByIds(context.Background(), &userpb.GetUsersRequest{UserIds: userIDs})
		if err != nil {
			log.Fatalf("Error while calling GetUsersByIds: %v", err)
		}

		// Process the response and send it to the client
		for {
			user, err := getUsersResponse.Recv()
			if err != nil {
				break
			}
			fmt.Fprintf(w, "User Details: %+v\n", user)
		}
	})

	// Modify the HTTP handler to accept JSON data
	mux.HandleFunc("/CreateUser", func(w http.ResponseWriter, r *http.Request) {
		// Read the JSON data from the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Error reading request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Decode the JSON data into a UserJSON struct
		var userJSON userpb.User
		if err := json.Unmarshal(body, &userJSON); err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// Call the gRPC server to create or update the user using userJSON
		// Implement your logic to store the user data as needed
		user := &userpb.User{
			Id:      userJSON.Id,
			Fname:   userJSON.Fname,
			City:    userJSON.City,
			Phone:   userJSON.Phone,
			Height:  float32(userJSON.Height),
			Married: userJSON.Married,
		}

		client.CreateUser(context.Background(), user)

		// Implement your user creation or update logic and store user data

		// Respond to the client as needed
		fmt.Fprintf(w, "User Created/Updated: %+v\n", user)
	})
	fmt.Println("Server Started at port 3333")
	http.ListenAndServe(":3333", mux)
}
