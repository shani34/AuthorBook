package book

import (
	"context"
	"database/sql"
	"log"

	"projects/GoLang-Interns-2022/authorbook/entities"
)

type Store struct {
	DB *sql.DB
}

// New : factory function used for dependency injection
func New(db *sql.DB) Store {
	return Store{db}
}

// GetAllBook : fetches the all book from database
func (bs Store) GetAllBook(ctx context.Context) ([]entities.Book, error) {
	var (
		books []entities.Book
		rows  *sql.Rows
		err   error
	)

	rows, err = bs.DB.QueryContext(ctx, "SELECT * FROM book")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book entities.Book

		err = rows.Scan(&book.BookID, &book.AuthorID, &book.Title, &book.Publication, &book.PublishedDate)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

// GetBooksByTitle : give the books with particular title
func (bs Store) GetBooksByTitle(ctx context.Context, title string) ([]entities.Book, error) {
	var (
		books []entities.Book
		Rows  *sql.Rows
		err   error
	)

	Rows, err = bs.DB.QueryContext(ctx, "SELECT * FROM book WHERE title=?", title)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer Rows.Close()

	for Rows.Next() {
		var book entities.Book

		err = Rows.Scan(&book.BookID, &book.AuthorID, &book.Title, &book.Publication, &book.PublishedDate)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

// GetBookByID : give the book with particular id
func (bs Store) GetBookByID(ctx context.Context, id int) (entities.Book, error) {
	var book entities.Book

	row := bs.DB.QueryRowContext(ctx, "select * from book where id=?", id)

	err := row.Scan(&book.BookID, &book.AuthorID, &book.Title, &book.Publication, &book.PublishedDate)
	if err != nil {
		log.Print(err)
		return entities.Book{}, err
	}

	return book, nil
}

// Post : inserts the book into database
func (bs Store) Post(ctx context.Context, book *entities.Book) (int, error) {
	result, err := bs.DB.ExecContext(ctx, "insert into book(author_id,title,publication,published_date)values(?,?,?,?)",
		book.AuthorID, book.Title, book.Publication, book.PublishedDate)
	if err != nil {
		log.Print(err)
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return -1, err
	}

	return int(id), nil
}

// Put : updates the book with particular id
func (bs Store) Put(ctx context.Context, book *entities.Book, id int) (int, error) {
	res, err := bs.DB.ExecContext(ctx, "update book set author_id=?,title=?,publication=?,published_date=? where id=?",
		book.AuthorID, book.Title, book.Publication, book.PublishedDate, id)
	if err != nil {
		return 0, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(ra), nil
}

// Delete : deletes the book by particular id
func (bs Store) Delete(ctx context.Context, id int) (int, error) {
	res, err := bs.DB.ExecContext(ctx, "delete from book where id=?", id)
	if err != nil {
		return -1, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(ra), nil
}
