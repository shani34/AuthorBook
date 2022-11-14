package book

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"projects/GoLang-Interns-2022/authorbook/entities"
)

// TestGetAllBook : to test GetAllBook
func TestGetAllBook(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Print(err)
	}

	var (
		book1 = entities.Book{BookID: 1, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{},
		}

		book2 = entities.Book{BookID: 2, AuthorID: 1, Title: "book two", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{},
		}

		books = sqlmock.NewRows([]string{"id", "author_id", "title", "publication", "published_date"}).
			AddRow(book1.BookID, book1.AuthorID, book1.Title, book1.Publication, book1.PublishedDate).AddRow(book2.BookID,
			book2.AuthorID, book2.Title, book2.Publication, book2.PublishedDate)
	)

	Testcases := []struct {
		desc string

		expected    []entities.Book
		expectedErr error
	}{
		{desc: "getting all books", expected: []entities.Book{book1, book2}, expectedErr: nil},
		{desc: "getting all books", expected: []entities.Book{}, expectedErr: errors.New("syntax error")},
	}

	for _, tc := range Testcases {
		bs := New(db)

		mock.ExpectQuery("SELECT * FROM book").WithArgs().WillReturnRows(books).WillReturnError(tc.expectedErr)

		b, err := bs.GetAllBook(context.TODO())
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(b, tc.expected) {
			t.Errorf("failed for %s ", tc.desc)
		}
	}
}

// TestGetBooksByTitle : to test GetBooksByTitle
func TestGetBooksByTitle(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Print(err)
	}

	var (
		book1 = entities.Book{BookID: 1, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{},
		}

		book2 = entities.Book{BookID: 2, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{},
		}

		book3 = entities.Book{}

		books = sqlmock.NewRows([]string{"id", "author_id", "title", "publication", "published_date"}).
			AddRow(book1.BookID, book1.AuthorID, book1.Title, book1.Publication, book1.PublishedDate).AddRow(book2.BookID,
			book2.AuthorID, book2.Title, book2.Publication, book2.PublishedDate)

		books1 = sqlmock.NewRows([]string{"id", "author_id", "title", "publication", "published_date"}).
			AddRow(book3.BookID, book3.AuthorID, book3.Title, book3.Publication, book3.PublishedDate)
	)

	Testcases := []struct {
		desc  string
		title string

		expected    []entities.Book
		expectedErr error
	}{
		{desc: "getting all books", title: "book one", expected: []entities.Book{book1, book2}, expectedErr: nil},
		{desc: "invalid case", title: "", expected: []entities.Book{}, expectedErr: errors.New("syntax error")},
		{desc: "scan error", title: "unique", expected: []entities.Book{book3}, expectedErr: nil},
	}

	for _, tc := range Testcases {
		bs := New(db)

		if tc.title != "unique" {
			mock.ExpectQuery("SELECT * FROM book WHERE title=?").WithArgs().WillReturnRows(books).
				WillReturnError(tc.expectedErr)
		} else {
			mock.ExpectQuery("SELECT * FROM book WHERE title=?").WithArgs().WillReturnRows(books1).
				WillReturnError(tc.expectedErr)
		}

		b, err := bs.GetBooksByTitle(context.TODO(), tc.title)

		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(b, tc.expected) {
			t.Errorf("failed for %s", tc.desc)
		}
	}
}

// TestGetBookByID : to test GetBookByID
func TestGetBookByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Print(err)
	}

	var (
		book = entities.Book{BookID: 1,
			AuthorID: 1, Title: "book one", Publication: "penguin", PublishedDate: "20/06/2000"}
		book1 = sqlmock.NewRows([]string{"id", "author_id", "title", "publication", "published_date"}).
			AddRow(book.BookID, book.AuthorID, book.Title, book.Publication, book.PublishedDate)
	)

	Testcases := []struct {
		desc     string
		targetID int

		expected    entities.Book
		expectedErr error
	}{
		{desc: "fetching book by id",
			targetID: 1, expected: entities.Book{BookID: 1,
				AuthorID: 1, Title: "book one", Publication: "penguin", PublishedDate: "20/06/2000"}, expectedErr: nil},

		{"invalid id", -1, entities.Book{}, errors.New("invalid")},
	}

	for _, tc := range Testcases {
		bs := New(db)

		mock.ExpectQuery("select * from book where id=?").WithArgs(tc.targetID).WillReturnRows(book1).WillReturnError(tc.expectedErr)

		b, err := bs.GetBookByID(context.TODO(), tc.targetID)
		if err != nil {
			log.Print(err)
		}

		if b != tc.expected || err != tc.expectedErr {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestPost : to test the post
func TestPost(t *testing.T) {
	testcases := []struct {
		desc  string
		input entities.Book

		expectedErr  error
		RowAffected  int64
		LastInserted int64
	}{
		{desc: "valid book", input: entities.Book{BookID: 1, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{}},
			expectedErr: nil, RowAffected: 1, LastInserted: 15,
		},
		{desc: "exiting book", input: entities.Book{BookID: 1, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{}},
			expectedErr: errors.New("already exists"), RowAffected: 0, LastInserted: 0,
		},
		{desc: "error case", input: entities.Book{BookID: 3, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{}},
			expectedErr: errors.New("last inserted error"), RowAffected: 1, LastInserted: 15,
		},
	}

	for _, tc := range testcases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error during the opening of database:%v\n", err)
		}

		if tc.input.BookID != 3 {
			mock.ExpectExec("insert into book(author_id,title,publication,published_date)values(?,?,?,?)").
				WithArgs(tc.input.AuthorID, tc.input.Title, tc.input.Publication, tc.input.PublishedDate).
				WillReturnResult(sqlmock.NewResult(tc.LastInserted, tc.RowAffected)).WillReturnError(tc.expectedErr)
		} else {
			mock.ExpectExec("insert into book(author_id,title,publication,published_date)values(?,?,?,?)").
				WithArgs(tc.input.AuthorID, tc.input.Title, tc.input.Publication, tc.input.PublishedDate).
				WillReturnResult(sqlmock.NewErrorResult(tc.expectedErr)).WillReturnError(nil)
		}

		bs := New(db)

		_, err = bs.Post(context.TODO(), &tc.input)
		if err != tc.expectedErr {
			t.Errorf("failed for %s", tc.desc)
		}
	}
}

// TestPut : to test the put
func TestPut(t *testing.T) {
	testcases := []struct {
		desc     string
		input    entities.Book
		targetID int

		expectedErr  error
		RowAffected  int64
		LastInserted int64
	}{
		{desc: "not existing book", input: entities.Book{BookID: 1, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{}}, targetID: -1,
			expectedErr: errors.New("does not exist"), RowAffected: 0, LastInserted: 0,
		},
		{desc: "exiting book", input: entities.Book{BookID: 12, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{}}, targetID: 4,
			expectedErr: nil, RowAffected: 1, LastInserted: 15,
		},
		{desc: "error case", input: entities.Book{BookID: 13, AuthorID: 1, Title: "book one", Publication: "penguin",
			PublishedDate: "20/06/2000", Author: entities.Author{}}, targetID: 4,
			expectedErr: errors.New("database error"), RowAffected: 1, LastInserted: 15,
		},
	}

	for _, tc := range testcases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error during the opening of database:%v\n", err)
		}

		if tc.input.BookID != 13 {
			mock.ExpectExec("update book set id=?,author_id=?,title=?,publication=?,published_date=? where id=?").
				WithArgs(tc.input.BookID, tc.input.AuthorID, tc.input.Title, tc.input.Publication, tc.input.PublishedDate, tc.targetID).
				WillReturnResult(sqlmock.NewResult(tc.LastInserted, tc.RowAffected)).WillReturnError(tc.expectedErr)
		} else {
			mock.ExpectExec("update book set id=?,author_id=?,title=?,publication=?,published_date=? where id=?").
				WithArgs(tc.input.BookID, tc.input.AuthorID, tc.input.Title, tc.input.Publication, tc.input.PublishedDate, tc.targetID).
				WillReturnResult(sqlmock.NewErrorResult(tc.expectedErr)).WillReturnError(nil)
		}

		bs := New(db)

		_, err = bs.Put(context.TODO(), &tc.input, tc.targetID)
		if err != tc.expectedErr {
			t.Errorf("failed for %s", tc.desc)
		}
	}
}

// TestDelete : to test delete method
func TestDelete(t *testing.T) {
	testcases := []struct {
		// input
		desc   string
		target int
		// output
		rowsAffected   int64
		lastInsertedID int64
		expectedErr    error
	}{
		{"valid authorId", 4, 1, 0, nil},
		{"invalid authorId", -1, 0, 0, errors.New("invalid bookID")},
		{"not existing", 100, 0, 0, errors.New("does not exist")},
	}

	for _, tc := range testcases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Print(err)
		}

		bs := New(db)

		if tc.target != 100 {
			mock.ExpectExec("delete from book where id=?").WithArgs(tc.target).
				WillReturnResult(sqlmock.NewResult(tc.lastInsertedID, tc.rowsAffected)).WillReturnError(tc.expectedErr)
		} else {
			mock.ExpectExec("delete from book where id=?").WithArgs(tc.target).
				WillReturnResult(sqlmock.NewErrorResult(tc.expectedErr)).WillReturnError(nil)
		}

		_, err = bs.Delete(context.TODO(), tc.target)
		if err != tc.expectedErr {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}
