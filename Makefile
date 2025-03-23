DOCKER_IMAGE_NAME=go-grpc-test

docker-build:
	docker build -t ${DOCKER_IMAGE_NAME} .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative book.proto

test:
	go test ./server