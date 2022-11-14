package bookhttp

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"encoding/json"
	"io"
	"strconv"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/service"
)

type BookHandler struct {
	bookH service.BookService
}

// New : factory function
func New(bookS service.BookService) BookHandler {
	return BookHandler{bookS}
}

// GetAllBook : handles the request of getting all books
func (h BookHandler) GetAllBook(ctx *gofr.Context) (interface{}, error) {
	includeAuthor := ctx.Param("includeAuthor")
	title := ctx.Param("title")

	books, err := h.bookH.GetAllBook(ctx, title, includeAuthor)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// GetBookByID : handles the request of getting a book
func (h BookHandler) GetBookByID(ctx *gofr.Context) (interface{}, error) {
	params := ctx.PathParam("id")

	id, err := strconv.Atoi(params)
	if err != nil || id < 0 {
		return nil, err
	}

	book, err := h.bookH.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// Post : handles the request of posting a book
func (h BookHandler) Post(ctx *gofr.Context) (interface{}, error) {
	var book entities.Book

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &book)
	if err != nil {
		return nil, err
	}

	book1, err := h.bookH.Post(ctx, &book)
	if err != nil {
		return nil, err
	}

	return book1, nil
}

// Put : handle the request of updating a book
func (h BookHandler) Put(ctx *gofr.Context) (interface{}, error) {
	var (
		body []byte
		book entities.Book
		id   int
	)

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &book)
	if err != nil {
		return nil, err
	}

	params := ctx.PathParam("id")

	id, err = strconv.Atoi(params)
	if err != nil {
		return nil, err
	}

	book, err = h.bookH.Put(ctx, &book, id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// Delete : handles the request of removing a book
func (h BookHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	params := ctx.PathParam("id")

	id, err := strconv.Atoi(params)
	if err != nil {
		return nil, err
	}

	err = h.bookH.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "successfully deleted", nil
}
