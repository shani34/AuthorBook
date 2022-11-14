package store

import (
	"context"
	"projects/GoLang-Interns-2022/authorbook/entities"
)

type AuthorStorer interface {
	Post(ctx context.Context, author entities.Author) (int, error)
	Put(ctx context.Context, author entities.Author, id int) (int, error)
	Delete(ctx context.Context, id int) (int, error)
	IncludeAuthor(ctx context.Context, id int) (entities.Author, error)
}

type BookStorer interface {
	GetAllBook(ctx context.Context) ([]entities.Book, error)
	GetBooksByTitle(ctx context.Context, title string) ([]entities.Book, error)

	GetBookByID(ctx context.Context, id int) (entities.Book, error)
	Post(ctx context.Context, book *entities.Book) (int, error)
	Put(ctx context.Context, book *entities.Book, id int) (int, error)
	Delete(ctx context.Context, id int) (int, error)
}
