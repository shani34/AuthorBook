package authorhttp

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"encoding/json"
	"io"
	"log"
	"strconv"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/service"
)

type AuthorHandler struct {
	authorService service.AuthorService
}

// New : factory function used for injection
func New(a service.AuthorService) AuthorHandler {
	return AuthorHandler{a}
}

// Post : handles the request of posting an author
func (h AuthorHandler) Post(c *gofr.Context) (interface{}, error) {
	var author entities.Author

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		//c.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = json.Unmarshal(body, &author)
	if err != nil {
		return nil, err
	}

	a, err := h.authorService.Post(c, author)
	if err != nil {
		log.Print("3")
		return nil, err
	}
	log.Print("4")
	return a, nil

}

// Put : handles the request of updating an author
func (h AuthorHandler) Put(ctx *gofr.Context) (interface{}, error) {
	var author entities.Author

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &author)
	if err != nil {
		return nil, err
	}

	params := ctx.PathParam("id")

	id, err := strconv.Atoi(params)
	if err != nil {
		return nil, err
	}

	author1, err := h.authorService.Put(ctx, author, id)
	if err != nil {
		return nil, err
	}

	return author1, nil

}

// Delete : handles the request of deleting an author
func (h AuthorHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	err = h.authorService.Delete(ctx, intID)
	if err != nil {
		return nil, err
	}

	return "successfully deleted!", nil
}
