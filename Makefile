OUTPUT = main 
SERVICE_NAME = transactions

build-local:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)/main.go

test:
	go test ./...

run: build-local
	@echo ">> Running application ..."
	DB_PORT=3306 \
	DB_HOST=localhost \
	DB_USER=root \
	DB_PASSWORD=password \
	DB_NAME=core \
	APP_PORT=9000 \
	./$(OUTPUT)