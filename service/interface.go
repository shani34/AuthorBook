package service

import (
	"context"
	"projects/GoLang-Interns-2022/authorbook/entities"
)

type AuthorService interface {
	Post(ctx context.Context, author entities.Author) (entities.Author, error)
	Put(ctx context.Context, author entities.Author, id int) (entities.Author, error)
	Delete(ctx context.Context, id int) error
}

type BookService interface {
	GetAllBook(ctx context.Context, title string, includeAuthor string) ([]entities.Book, error)
	GetBookByID(ctx context.Context, id int) (entities.Book, error)
	Post(ctx context.Context, book *entities.Book) (entities.Book, error)
	Put(ctx context.Context, book *entities.Book, id int) (entities.Book, error)
	Delete(ctx context.Context, id int) error
}
