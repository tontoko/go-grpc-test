package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	book "go-grpc-test" // 生成されたコード
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := book.NewBookServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateBook(ctx, &book.CreateBookRequest{
		Title:           "Test Book",
		Author:          "Test Author",
		Isbn:            "978-0321765723",
		PublicationDate: "2024-01-01",
		Genre:           "Test Genre",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Book: %s", r.GetTitle())

	// Get Book
	getBookResponse, err := c.GetBook(ctx, &book.GetBookRequest{Isbn: "978-0321765723"})
	if err != nil {
		log.Fatalf("could not get book: %v", err)
	}
	log.Printf("Get Book: %s", getBookResponse.GetTitle())

	// Update Book
	updateBookResponse, err := c.UpdateBook(ctx, &book.UpdateBookRequest{
		Isbn:            "978-0321765723",
		Title:           "Updated Test Book",
		Author:          "Updated Test Author",
		PublicationDate: "2024-02-02",
		Genre:           "Updated Test Genre",
	})
	if err != nil {
		log.Fatalf("could not update book: %v", err)
	}
	log.Printf("Update Book: %s", updateBookResponse.GetTitle())

	// Delete Book
	_, err = c.DeleteBook(ctx, &book.DeleteBookRequest{Isbn: "978-0321765723"})
	if err != nil {
		log.Fatalf("could not delete book: %v", err)
	}
	log.Printf("Book Deleted")

	os.Exit(0)
}
