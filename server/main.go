package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	book "go-grpc-test" // 生成されたコード
)

const (
	grpcPort = ":50051"
	httpPort = ":8080"
)

type BookServer struct {
	book.UnimplementedBookServiceServer
	db *gorm.DB
}

type Book struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	Isbn            string `json:"isbn" gorm:"primaryKey"`
	PublicationDate string `json:"publication_date"`
	Genre           string `json:"genre"`
}

func (s *BookServer) CreateBook(ctx context.Context, req *book.CreateBookRequest) (*book.Book, error) {
	b := Book{
		Title:           req.Title,
		Author:          req.Author,
		Isbn:            req.Isbn,
		PublicationDate: req.PublicationDate,
		Genre:           req.Genre,
	}
	result := s.db.Create(&b)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to create book: %v", result.Error)
	}
	return &book.Book{
		Title:           b.Title,
		Author:          b.Author,
		Isbn:            b.Isbn,
		PublicationDate: b.PublicationDate,
		Genre:           b.Genre,
	}, nil
}

func (s *BookServer) GetBook(ctx context.Context, req *book.GetBookRequest) (*book.Book, error) {
	var b Book
	result := s.db.First(&b, "isbn = ?", req.Isbn)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "book not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get book: %v", result.Error)
	}
	return &book.Book{
		Title:           b.Title,
		Author:          b.Author,
		Isbn:            b.Isbn,
		PublicationDate: b.PublicationDate,
		Genre:           b.Genre,
	}, nil
}

func (s *BookServer) UpdateBook(ctx context.Context, req *book.UpdateBookRequest) (*book.Book, error) {
	var b Book
	result := s.db.First(&b, "isbn = ?", req.Isbn)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "book not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get book: %v", result.Error)
	}

	b.Title = req.Title
	b.Author = req.Author
	b.PublicationDate = req.PublicationDate
	b.Genre = req.Genre

	result = s.db.Save(&b)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update book: %v", result.Error)
	}

	return &book.Book{
		Title:           b.Title,
		Author:          b.Author,
		Isbn:            b.Isbn,
		PublicationDate: b.PublicationDate,
		Genre:           b.Genre,
	}, nil
}

func (s *BookServer) DeleteBook(ctx context.Context, req *book.DeleteBookRequest) (*book.Empty, error) {
	result := s.db.Delete(&Book{}, "isbn = ?", req.Isbn)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete book: %v", result.Error)
	}
	return &book.Empty{}, nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("book.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// マイグレーション
	db.AutoMigrate(&Book{})

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bookServer := &BookServer{db: db}
	book.RegisterBookServiceServer(s, bookServer)

	log.Printf("gRPC server listening on %s", grpcPort)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := gin.Default()

	r.POST("/books", func(c *gin.Context) {
		log.Println("Received POST /books request")
		var req Book
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		conn, err := grpc.Dial("localhost"+grpcPort, grpc.WithInsecure())
		log.Printf("Attempting to connect to gRPC server at localhost%s", grpcPort)
		if err != nil {
			log.Printf("did not connect: %v", err)
			c.JSON(500, gin.H{"error": "grpc server not running"})
			return
		}
		defer conn.Close()

		client := book.NewBookServiceClient(conn)
		grpcReq := &book.CreateBookRequest{Title: req.Title, Author: req.Author, Isbn: req.Isbn, PublicationDate: req.PublicationDate, Genre: req.Genre}

		log.Printf("Sending CreateBook request to gRPC server: %v", grpcReq)
		res, err := client.CreateBook(context.Background(), grpcReq)
		if err != nil {
			log.Printf("could not greet: %v", err)

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, res)
		log.Println("Successfully created book via gRPC")
	})

	r.GET("/books/:isbn", func(c *gin.Context) {
		isbn := c.Param("isbn")
		conn, err := grpc.Dial("localhost"+grpcPort, grpc.WithInsecure())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()

		client := book.NewBookServiceClient(conn)
		req := &book.GetBookRequest{Isbn: isbn}
		res, err := client.GetBook(context.Background(), req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	log.Printf("HTTP server listening on %s", httpPort)
	if err := r.Run("0.0.0.0" + httpPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
