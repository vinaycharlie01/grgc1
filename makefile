

run:
	protoc -Iuser/proto --go_opt=module=myapp --go_out=. --go-grpc_opt=module=myapp --go-grpc_out=. user/proto/*.proto
	go build -o bin/greet/server ./greet/server
	go build -o bin/greet/client ./greet/client