package main

import (
	books "GRPC_API_Go_Books/Books"
	"context"
	"fmt"
	"log"
	"strconv"

	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/cobra"
)

func printTimeCmd() *cobra.Command {
	return &cobra.Command{
		Use: "curtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now()
			prettyTime := now.Format(time.RubyDate)
			cmd.Println("Hey! The current time is ", prettyTime)
			return nil //Tells cobra that no errors happened
		},
	}
}

func getBooksCmd() *cobra.Command {
	return &cobra.Command{
		Use: "books",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *grpc.ClientConn
			conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()

			b := books.NewBooksClient(conn)
			request := books.GetBooksRequest{}

			response, err := b.GetBooks(context.Background(), &request)
			if err != nil {
				log.Fatalf("Error when calling GetBooks: %s", err)
			}

			for _, book := range response.Book {
				fmt.Println(book)
			}
			return nil //Tells cobra that no errors happened
		},
	}
}

func getBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "book",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *grpc.ClientConn
			conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()

			b := books.NewBooksClient(conn)
			book_id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatalf("Error when calling GetBook, invalid ld: %s", err)
			}
			request := books.GetBookRequest{Id: book_id}
			response, err := b.GetBook(context.Background(), &request)
			if err != nil {
				log.Fatalf("Error when calling GetBook: %s", err)
			}
			fmt.Println(response.Book)
			return nil //Tells cobra that no errors happened
		},
	}
}

func addBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "addbook",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *grpc.ClientConn
			conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()

			b := books.NewBooksClient(conn)
			newBook := books.Book{Id: 4, Title: "Colors for adults", Author: "Adeline", Quantity: 2}
			request := books.CreateBookRequest{Book: &newBook}
			response, err := b.CreateBook(context.Background(), &request)
			if err != nil {
				log.Fatalf("Error when calling CreateBook: %s", err)
			}
			fmt.Printf("New book added successfuly: ")
			fmt.Println(response.Book)
			return nil //Tells cobra that no errors happened
		},
	}
}

func deleteBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "deletebook",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *grpc.ClientConn
			conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()

			b := books.NewBooksClient(conn)
			book_id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatalf("Error when calling DeleteBook, invalid ID: %s", err)
			}
			request := books.DeleteBookRequest{Id: book_id}
			response, err := b.DeleteBook(context.Background(), &request)
			if err != nil {
				log.Fatalf("Error when calling DeleteBook: %s", err)
			}
			fmt.Printf("You deleted a book succesfully: ")
			fmt.Println(response.Book)
			return nil //Tells cobra that no errors happened
		},
	}
}

func checkoutBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "checkout",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *grpc.ClientConn
			conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()

			b := books.NewBooksClient(conn)
			book_id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatalf("Error when calling CheckoutBook, invalid ld: %s", err)
			}
			request := books.CheckoutBookRequest{Id: book_id}
			response, err := b.CheckoutBook(context.Background(), &request)
			if err != nil {
				log.Fatalf("Error when calling CheckoutBook: %s", err)
			}
			fmt.Printf("You checked out a book succesfully: ")
			fmt.Println(response.Book)
			return nil //Tells cobra that no errors happened
		},
	}
}

func returnBookCmd() *cobra.Command {
	return &cobra.Command{
		Use: "return",
		RunE: func(cmd *cobra.Command, args []string) error {
			var conn *grpc.ClientConn
			conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()

			b := books.NewBooksClient(conn)
			book_id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatalf("Error when calling ReturnBook, invalid ld: %s", err)
			}
			request := books.ReturnBookRequest{Id: book_id}
			response, err := b.ReturnBook(context.Background(), &request)
			if err != nil {
				log.Fatalf("Error when calling ReturnBook: %s", err)
			}
			fmt.Printf("You returned a book succesfully: ")
			fmt.Println(response.Book)
			return nil //Tells cobra that no errors happened
		},
	}
}

func main() {
	cmd := &cobra.Command{
		Use:   "gifm",
		Short: "Welcome to the Client!",
	}

	cmd.AddCommand(printTimeCmd())
	cmd.AddCommand(getBooksCmd())
	cmd.AddCommand(getBookCmd())
	cmd.AddCommand(addBookCmd())
	cmd.AddCommand(deleteBookCmd())
	cmd.AddCommand(checkoutBookCmd())
	cmd.AddCommand(returnBookCmd())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
