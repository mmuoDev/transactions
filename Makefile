OUTPUT = main 
SERVICE_NAME = transactions

build-local:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)/main.go

test:
	go test ./...

clean:
	rm -f $(OUTPUT)

run: build-local
	@echo ">> Running application ..."
	DB_PORT= \
	DB_HOST= \
	DB_USER= \
	DB_PASSWORD= \
	DB_NAME= \
	APP_PORT=9000 \
	./$(OUTPUT)