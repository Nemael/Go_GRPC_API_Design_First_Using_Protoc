package main

import (
	books "GRPC_API_Go_Books/Books"
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
)

type myBooksServer struct {
	books.UnimplementedBooksServer
}

var (
	ConnectionString = "rest:password@tcp(localhost:3306)/books"
)

func getDB() (*sql.DB, error) {
	return sql.Open("mysql", ConnectionString)
}

func (s myBooksServer) GetBooks(ctx context.Context, req *books.GetBooksRequest) (*books.GetBooksResponse, error) {
	my_books := []*books.Book{}
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, title, author, quantity FROM books")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var book books.Book
		rows.Scan(&book.Id, &book.Title, &book.Author, &book.Quantity)
		my_books = append(my_books, &book)
	}
	return &books.GetBooksResponse{
		Book: my_books,
	}, nil
}

func (s myBooksServer) GetBook(ctx context.Context, req *books.GetBookRequest) (*books.GetBookResponse, error) {
	var my_book books.Book
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT id, title, author, quantity FROM books WHERE id = ?", int64(req.Id))
	err = row.Scan(&my_book.Id, &my_book.Title, &my_book.Author, &my_book.Quantity)
	if err != nil {
		return nil, err
	}
	return &books.GetBookResponse{
		Book: &my_book,
	}, nil
}

func (s myBooksServer) CreateBook(ctx context.Context, req *books.CreateBookRequest) (*books.CreateBookResponse, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	newBook := books.Book{Id: req.Book.Id, Title: req.Book.Title, Author: req.Book.Author, Quantity: req.Book.Quantity}
	_, err = db.Exec("INSERT INTO books (id, title, author, quantity) VALUES (?, ?, ?, ?)", newBook.Id, newBook.Title, newBook.Author, newBook.Quantity)
	if err != nil {
		return nil, err
	}
	return &books.CreateBookResponse{
		Book: &newBook,
	}, nil
}

func (s myBooksServer) CheckoutBook(ctx context.Context, req *books.CheckoutBookRequest) (*books.CheckoutBookResponse, error) {
	var book books.Book
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	id := req.Id // Query parameter
	row := db.QueryRow("SELECT id, title, author, quantity FROM books where id = ?", id)
	err = row.Scan(&book.Id, &book.Title, &book.Author, &book.Quantity)
	if err != nil {
		return nil, err
	}

	if book.Quantity <= int64(0) {
		return nil, errors.New("500: No copies of this book available")
	}
	book.Quantity -= int64(1)

	db.QueryRow("UPDATE books SET quantity = ? WHERE id = ?", book.Quantity, book.Id)
	return &books.CheckoutBookResponse{
		Book: &book,
	}, nil
}

func (s myBooksServer) ReturnBook(ctx context.Context, req *books.ReturnBookRequest) (*books.ReturnBookResponse, error) {
	var book books.Book
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	id := req.Id // Query parameter
	row := db.QueryRow("SELECT id, title, author, quantity FROM books where id = ?", id)
	err = row.Scan(&book.Id, &book.Title, &book.Author, &book.Quantity)
	if err != nil {
		return nil, err
	}
	book.Quantity += int64(1)

	db.QueryRow("UPDATE books SET quantity = ? WHERE id = ?", book.Quantity, book.Id)
	return &books.ReturnBookResponse{
		Book: &book,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Cannot create listener %s", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myBooksServer{}

	books.RegisterBooksServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
