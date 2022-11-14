package author

import (
	"context"
	"database/sql"
	"log"

	"projects/GoLang-Interns-2022/authorbook/entities"
)

type Store struct {
	DB *sql.DB
}

// New : factory function
func New(db *sql.DB) Store {
	return Store{db}
}

// Post : insert an author
func (s Store) Post(ctx context.Context, author entities.Author) (int, error) {
	res, err := s.DB.ExecContext(ctx, "insert into author(first_name,last_name,dob,pen_name)values(?,?,?,?)",
		author.FirstName, author.LastName, author.DOB, author.PenName)
	if err != nil {
		log.Print(err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

// Put : inserts an author if that does not exist and update author if exists
func (s Store) Put(ctx context.Context, author entities.Author, id int) (int, error) {
	_, err := s.DB.ExecContext(ctx, "update author set first_name=?,last_name=?,dob=?,pen_name=? where author_id=?",
		author.FirstName, author.LastName, author.DOB, author.PenName, id)
	if err != nil {
		log.Print(err)
		return -1, err
	}

	return id, nil
}

// Delete :  deletes an author
func (s Store) Delete(ctx context.Context, id int) (int, error) {
	res, err := s.DB.ExecContext(ctx, "delete from author where author_id=?", id)
	if err != nil {
		return -1, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return int(ra), nil
}

// IncludeAuthor : checks whether an author exists or not if exists then it returns the author detail
func (s Store) IncludeAuthor(ctx context.Context, id int) (entities.Author, error) {
	var author entities.Author

	Row := s.DB.QueryRowContext(ctx, "SELECT * FROM author where author_id=?", id)

	if err := Row.Scan(&author.AuthorID, &author.FirstName, &author.LastName, &author.DOB, &author.PenName); err != nil {
		return entities.Author{}, err
	}

	return author, nil
}
