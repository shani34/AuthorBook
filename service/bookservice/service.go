package bookservice

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/store"
)

type BookService struct {
	bookService   store.BookStorer
	authorService store.AuthorStorer
}

// New : factory function
func New(bs store.BookStorer, as store.AuthorStorer) BookService {
	return BookService{bs, as}
}

// GetAllBook : implements the logic of getting all book
func (b BookService) GetAllBook(ctx context.Context, title, includeAuthor string) ([]entities.Book, error) {
	var (
		books []entities.Book
		err   error
	)

	if title != "" {
		books, err = b.bookService.GetBooksByTitle(ctx, title)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	} else {
		books, err = b.bookService.GetAllBook(ctx)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	if includeAuthor == "true" {
		for i := range books {
			author, err := b.authorService.IncludeAuthor(ctx, books[i].AuthorID)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			books[i].Author = &author
		}
	}

	return books, nil
}

// GetBookByID : implements the logic of getting a single by
func (b BookService) GetBookByID(ctx context.Context, id int) (entities.Book, error) {
	if id <= 0 {
		return entities.Book{}, errors.New("invalid id")
	}

	book, err := b.bookService.GetBookByID(ctx, id)
	if err != nil {
		log.Print(err)
		return entities.Book{}, err
	}

	author, err := b.authorService.IncludeAuthor(ctx, book.AuthorID)
	book.Author = &author

	return book, nil
}

// Post : checks the book before posting
func (b BookService) Post(ctx context.Context, book *entities.Book) (entities.Book, error) {
	if book.Title == "" || book.AuthorID < 0 || checkPublication(book.Publication) ||
		!checkPublishedDate(book.PublishedDate) {
		return entities.Book{}, errors.New("invalid constraints")
	}

	existAuthor, err := b.authorService.IncludeAuthor(ctx, book.AuthorID)
	if err != nil {
		return entities.Book{}, err
	}

	id, err := b.bookService.Post(ctx, book)
	if err != nil || id == -1 {
		return entities.Book{}, errors.New("database issue")
	}

	book.Author = &existAuthor
	book.BookID = id

	return *book, nil
}

// Put :  checks the book before updating
func (b BookService) Put(ctx context.Context, book *entities.Book, id int) (entities.Book, error) {
	if book.Title == "" || book.AuthorID <= 0 || checkPublication(book.Publication) ||
		!checkPublishedDate(book.PublishedDate) {
		return entities.Book{}, errors.New("invalid constraints")
	}

	author, err := b.authorService.IncludeAuthor(ctx, book.AuthorID)
	if err != nil {
		return entities.Book{}, errors.New("author does not exist")
	}

	count, err := b.bookService.Put(ctx, book, id)
	if err != nil || count <= 0 {
		return entities.Book{}, errors.New("book does not exist")
	}

	book.Author = &author
	book.BookID = id

	return *book, nil
}

// Delete : checks before deleting a book
func (b BookService) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return errors.New("invalid id")
	}

	count, err := b.bookService.Delete(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return errors.New("book does not exist")
	}

	return nil
}

// checkPublication : validates publication
func checkPublication(publication string) bool {
	publication = strings.ToLower(publication)

	return !(publication == "penguin" || publication == "scholastic" || publication == "arihant")
}

// checkPublishedDate : validates the published date
func checkPublishedDate(publishedDate string) bool {
	Dob := strings.Split(publishedDate, "/")
	day, _ := strconv.Atoi(Dob[0])
	month, _ := strconv.Atoi(Dob[1])
	year, _ := strconv.Atoi(Dob[2])

	switch {
	case day <= 0 || day > 31:
		return false
	case month <= 0 || month > 12:
		return false
	case year < 1870 || year > 2022:
		return false
	}

	return true
}
