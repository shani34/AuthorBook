package authorservice

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/store"
)

type AuthorService struct {
	datastore store.AuthorStorer
}

// New : factory function , use for dependency injection
func New(s store.AuthorStorer) AuthorService {
	return AuthorService{s}
}

// Post : checks the author before posting
func (s AuthorService) Post(ctx context.Context, a entities.Author) (entities.Author, error) {
	if a.FirstName == "" || !checkDob(a.DOB) {
		return entities.Author{}, errors.New("invalid constraints")
	}

	id, err := s.datastore.Post(ctx, a)
	if err != nil || id <= 0 {
		return entities.Author{}, err
	}

	a.AuthorID = id

	return a, nil
}

// Put : checks the author before updating
func (s AuthorService) Put(ctx context.Context, a entities.Author, id int) (entities.Author, error) {
	if a.FirstName == "" || !checkDob(a.DOB) {
		return entities.Author{}, errors.New("invalid constraints")
	}

	existAuthor, err := s.datastore.IncludeAuthor(ctx, id)
	if err != nil || existAuthor.AuthorID != id {
		return entities.Author{}, errors.New("author does not exist")
	}

	i, err := s.datastore.Put(ctx, a, id)
	if err != nil {
		return entities.Author{}, err
	}

	a.AuthorID = i

	return a, nil
}

// Delete : Deletes the author at particular id
func (s AuthorService) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return errors.New("invalid id")
	}

	count, err := s.datastore.Delete(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return errors.New("author does not exist")
	}

	return nil
}

// checkDob : validates the DOB
func checkDob(dob string) bool {
	Dob := strings.Split(dob, "/")
	day, _ := strconv.Atoi(Dob[0])
	month, _ := strconv.Atoi(Dob[1])
	year, _ := strconv.Atoi(Dob[2])

	switch {
	case day <= 0 || day > 31:
		return false
	case month <= 0 || month > 12:
		return false
	case year <= 0:
		return false
	}

	return true
}
