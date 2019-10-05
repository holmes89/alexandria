package main

import (
	"context"
	"time"
)

type Book struct {
	ID string
	DisplayName string
	Name string
	Path string
	Type string
	Description string
	Created time.Time
	Modified time.Time
}

type BookService interface {
	GetAll(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id string) (Book, error)
	Add(ctx context.Context, book *Book) error
}

type BookRepository interface{
	FindAll(ctx context.Context) ([]Book, error)
	FindByID(ctx context.Context, id string) (Book, error)
	Insert(ctx context.Context, book *Book) error
}