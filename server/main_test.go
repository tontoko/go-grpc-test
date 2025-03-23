package main

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	book "go-grpc-test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

func newTestServer(db *gorm.DB) *BookServer {
	s := &BookServer{db: db}
	return s
}

func newTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// マイグレーション
	db.AutoMigrate(&Book{})
	return db
}

func dialer(db *gorm.DB) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(bufSize)

	server := grpc.NewServer()
	bookServer := newTestServer(db)
	book.RegisterBookServiceServer(server, bookServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestBookService(t *testing.T) {
 	ctx := context.Background()
 	db := newTestDB()
 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()),
 		grpc.WithContextDialer(dialer(db)))
 	if err != nil {
 		t.Fatalf("Failed to dial bufnet: %v", err)
 	}
 	defer conn.Close()
 	client := book.NewBookServiceClient(conn)
 

 	// CreateBook
 	createRes, err := client.CreateBook(ctx, &book.CreateBookRequest{
 		Title:         "Test Book",
 		Author:        "Test Author",
 		Isbn:          "test-isbn",
 		PublicationDate: "2024-01-01",
 		Genre:         "Test Genre",
 	})
 	if err != nil {
 		t.Fatalf("CreateBook failed: %v", err)
 	}
 	if createRes.Title != "Test Book" {
 		t.Fatalf("CreateBook response title incorrect: %v", createRes.Title)
 	}
 

 	// GetBook
 	getRes, err := client.GetBook(ctx, &book.GetBookRequest{Isbn: "test-isbn"})
 	if err != nil {
 		t.Fatalf("GetBook failed: %v", err)
 	}
 	if getRes.Title != "Test Book" {
 		t.Fatalf("GetBook response title incorrect: %v", getRes.Title)
 	}
 

 	// UpdateBook
 	updateRes, err := client.UpdateBook(ctx, &book.UpdateBookRequest{
 		Title:         "Updated Book",
 		Author:        "Updated Author",
 		Isbn:          "test-isbn",
 		PublicationDate: "2024-02-01",
 		Genre:         "Updated Genre",
 	})
 	if err != nil {
 		t.Fatalf("UpdateBook failed: %v", err)
 	}
 	if updateRes.Title != "Updated Book" {
 		t.Fatalf("UpdateBook response title incorrect: %v", updateRes.Title)
 	}
 

 	// DeleteBook
 	_, err = client.DeleteBook(ctx, &book.DeleteBookRequest{Isbn: "test-isbn"})
 	if err != nil {
 		t.Fatalf("DeleteBook failed: %v", err)
 	}
 

 	// GetBook (確認)
 	_, err = client.GetBook(ctx, &book.GetBookRequest{Isbn: "test-isbn"})
 	if err == nil {
 		t.Fatalf("GetBook should have failed after deletion")
 	}
}

func TestMain(m *testing.M) {
 	// テスト実行前の処理 (必要に応じて)
 	exitCode := m.Run()
 	// テスト実行後の処理 (必要に応じて)
 	os.Exit(exitCode)
}