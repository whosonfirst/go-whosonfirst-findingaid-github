cli:
	go build -mod vendor -o bin/populate cmd/populate/main.go
	go build -mod vendor -o bin/update cmd/update/main.go
