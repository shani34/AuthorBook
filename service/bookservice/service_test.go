package bookservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/store"

	"github.com/golang/mock/gomock"
)

// TestGetAllBook : test the business logic of getting all book
func TestGetAllBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthorStore := store.NewMockAuthorStorer(ctrl)
	mockBookStore := store.NewMockBookStorer(ctrl)
	mock := New(mockBookStore, mockAuthorStore)

	Testcases := []struct {
		desc          string
		title         string
		includeAuthor string

		expected    []entities.Book
		expectedErr error
	}{
		{desc: "getting all books", title: "", includeAuthor: "", expected: []entities.Book{},
			expectedErr: errors.New("empty")},
		{desc: "getting book with author and particular", title: "book two", includeAuthor: "",
			expected: []entities.Book{}, expectedErr: errors.New("empty"),
		},
		{desc: "getting book with author and particular title", title: "book", includeAuthor: "true",
			expected: []entities.Book{}, expectedErr: errors.New("empty"),
		},
	}

	for _, tc := range Testcases {
		mockBookStore.EXPECT().GetAllBook(context.TODO()).Return(tc.expected, tc.expectedErr).AnyTimes()
		mockBookStore.EXPECT().GetBooksByTitle(context.TODO(), tc.title).Return(tc.expected, tc.expectedErr).AnyTimes()

		books, _ := mock.GetAllBook(context.TODO(), tc.title, tc.includeAuthor)

		if !reflect.DeepEqual(books, tc.expected) {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestGetBookByID : to test getBookByID
func TestGetBookByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthorStore := store.NewMockAuthorStorer(ctrl)
	mockBookStore := store.NewMockBookStorer(ctrl)
	mock := New(mockBookStore, mockAuthorStore)

	Testcases := []struct {
		desc     string
		targetID int

		expectedBody entities.Book
		expectedErr  error
	}{
		{desc: "fetching book by id",
			targetID: 1, expectedBody: entities.Book{}, expectedErr: errors.New("invalid id"),
		},
		{"invalid id", -1, entities.Book{}, errors.New("invalid id")},
	}

	for _, tc := range Testcases {
		mockBookStore.EXPECT().GetBookByID(context.TODO(), tc.targetID).Return(tc.expectedBody, tc.expectedErr).AnyTimes()
		book, _ := mock.GetBookByID(context.TODO(), tc.targetID)

		if !reflect.DeepEqual(book, tc.expectedBody) {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestPost : to test post method
func TestPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthorStore := store.NewMockAuthorStorer(ctrl)
	mockBookStore := store.NewMockBookStorer(ctrl)
	mock := New(mockBookStore, mockAuthorStore)

	testcases := []struct {
		desc  string
		input entities.Book

		expected     entities.Book
		expectedErr  error
		expectedErr1 error
	}{
		{desc: "success case", input: entities.Book{BookID: 0, AuthorID: 1, Title: "deciding decade",
			Publication: "penguin", PublishedDate: "20/03/2010", Author: &entities.Author{AuthorID: 1, FirstName: "shani",
				LastName: "kumar", DOB: "30/05/1999", PenName: "sk"}},
			expected: entities.Book{BookID: 12, AuthorID: 1, Title: "deciding decade", Publication: "penguin",
				PublishedDate: "20/03/2010", Author: &entities.Author{AuthorID: 1, FirstName: "shani",
					LastName: "kumar", DOB: "30/05/1999", PenName: "sk"}}, expectedErr: nil, expectedErr1: nil,
		},

		{desc: "author does not exist", input: entities.Book{BookID: 1, AuthorID: 3, Title: "deciding decade",
			Publication: "penguin", PublishedDate: "20/03/2010", Author: &entities.Author{}},
			expected: entities.Book{}, expectedErr: errors.New("issue"), expectedErr1: errors.New("author does not exist"),
		},

		{desc: "invalid publication", input: entities.Book{BookID: 1, AuthorID: 3, Title: "deciding decade",
			Publication: "pen", PublishedDate: "20/03/2010", Author: &entities.Author{}},
			expected: entities.Book{}, expectedErr: nil, expectedErr1: nil,
		},
	}
	for _, tc := range testcases {
		mockBookStore.EXPECT().Post(context.TODO(), &tc.input).Return(tc.expected.BookID, tc.expectedErr).AnyTimes()
		mockAuthorStore.EXPECT().IncludeAuthor(context.TODO(), tc.input.AuthorID).Return(tc.input.Author, tc.expectedErr1).AnyTimes()

		book, _ := mock.Post(context.TODO(), &tc.input)
		if !reflect.DeepEqual(book, tc.expected) {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestPut : to test the put method
func TestPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthorStore := store.NewMockAuthorStorer(ctrl)
	mockBookStore := store.NewMockBookStorer(ctrl)
	mock := New(mockBookStore, mockAuthorStore)

	testcases := []struct {
		desc    string
		input   entities.Book
		inputID int

		expected    entities.Book
		expectedErr error
	}{
		{desc: "success case", input: entities.Book{BookID: 12, AuthorID: 1, Title: "deciding decade",
			Publication: "penguin", PublishedDate: "20/03/2010", Author: &entities.Author{}}, inputID: 1,
			expected: entities.Book{BookID: 12, AuthorID: 1, Title: "deciding decade", Publication: "penguin",
				PublishedDate: "20/03/2010", Author: &entities.Author{}}, expectedErr: nil,
		},
		{desc: "invalid publication", input: entities.Book{BookID: 1, AuthorID: 1, Title: "deciding decade",
			Publication: "pen", PublishedDate: "20/03/2010", Author: &entities.Author{}},
			expected: entities.Book{}, expectedErr: nil,
		},
		{desc: "error", input: entities.Book{AuthorID: 1, Title: "deciding decade",
			Publication: "penguin", PublishedDate: "20/03/2010", Author: &entities.Author{}},
			expectedErr: errors.New("something went wrong"),
		},
	}
	for _, tc := range testcases {
		if tc.desc != "invalid publication" {
			mockBookStore.EXPECT().Put(context.TODO(), &tc.input, tc.inputID).Return(tc.expected.BookID, tc.expectedErr)
		}

		book, _ := mock.Put(context.TODO(), &tc.input, tc.inputID)

		if !reflect.DeepEqual(book, tc.expected) {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestDelete : to test delete method
func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthorStore := store.NewMockAuthorStorer(ctrl)
	mockBookStore := store.NewMockBookStorer(ctrl)
	mock := New(mockBookStore, mockAuthorStore)

	testcases := []struct {
		desc    string
		inputID int

		expectedID  int
		expectedErr error
	}{
		{"valid id", 1, 1, nil},
		{"invalid id", -1, -1, nil},
		{"error case", 1, -1, errors.New("something went wrong")},
	}

	for _, tc := range testcases {
		if tc.desc != "invalid id" {
			mockBookStore.EXPECT().Delete(context.TODO(), tc.inputID).Return(tc.expectedID, tc.expectedErr)
		}

		err := mock.Delete(context.TODO(), tc.inputID)

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}
