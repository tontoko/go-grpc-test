FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

CMD ["go", "run", "server/main.go"]