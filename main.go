package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	_ "github.com/go-sql-driver/mysql"
	"projects/GoLang-Interns-2022/authorbook/driver"
	"projects/GoLang-Interns-2022/authorbook/http/authorhttp"
	"projects/GoLang-Interns-2022/authorbook/http/bookhttp"
	"projects/GoLang-Interns-2022/authorbook/service/authorservice"
	"projects/GoLang-Interns-2022/authorbook/service/bookservice"
	"projects/GoLang-Interns-2022/authorbook/store/author"
	"projects/GoLang-Interns-2022/authorbook/store/book"
)

func main() {
	//r := mux.NewRouter()

	DB := driver.Connection()
	defer DB.Close()

	app := gofr.New()

	authorStore := author.New(DB)
	authorService := authorservice.New(authorStore)
	authorHandler := authorhttp.New(authorService)
	// author endpoints
	app.POST("/author", authorHandler.Post)
	app.DELETE("/author/{id}", authorHandler.Delete)
	app.PUT("/author/{id}", authorHandler.Put)

	bookStore := book.New(DB)
	bookService := bookservice.New(bookStore, authorStore)
	bookHandler := bookhttp.New(bookService)
	//book  endpoints
	app.GET("/book", bookHandler.GetAllBook)
	app.GET("/book/{id}", bookHandler.GetBookByID)
	app.POST("/book", bookHandler.Post)
	app.PUT("/book/{id}", bookHandler.Put)
	app.DELETE("/book/{id}", bookHandler.Delete)

	app.Start()
}
