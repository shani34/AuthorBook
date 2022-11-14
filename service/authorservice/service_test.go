package authorservice

import (
	"context"
	"errors"
	"testing"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/store"

	"github.com/golang/mock/gomock"
)

// TestPost : test the logic of posting an author
func TestPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := store.NewMockAuthorStorer(ctrl)
	mock := New(mockStore) // defining the type of interface

	testcases := []struct {
		desc string
		body entities.Author

		expectedAuthor entities.Author
		expectedID     int
		expectedErr    error
	}{
		{desc: "valid author", body: entities.Author{
			AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/05/1990", PenName: "Dark horse"},
			expectedAuthor: entities.Author{AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/05/1990",
				PenName: "Dark horse"}, expectedID: 4, expectedErr: nil},

		{desc: "existing author", body: entities.Author{
			AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal", DOB: "01/05/1990", PenName: "Dark horse"},
			expectedAuthor: entities.Author{}, expectedID: -1, expectedErr: errors.New("already exists")},

		{desc: "invalid firstname", body: entities.Author{
			AuthorID: 5, FirstName: "", LastName: "mrinal", DOB: "20/01/1990", PenName: "Dark horse"},
			expectedAuthor: entities.Author{}, expectedID: -1, expectedErr: errors.New("invalid constraints")},

		{desc: "invalid DOB", body: entities.Author{
			AuthorID: 5, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/01/0", PenName: "Dark horse"},
			expectedAuthor: entities.Author{}, expectedID: -1, expectedErr: errors.New("invalid constraints")},

		{desc: "invalid day", body: entities.Author{
			AuthorID: 5, FirstName: "nilotpal", LastName: "mrinal", DOB: "0/01/2000", PenName: "Dark horse"},
			expectedAuthor: entities.Author{}, expectedID: -1, expectedErr: errors.New("invalid constraints")},
	}

	for _, tc := range testcases {
		mockStore.EXPECT().Post(context.TODO(), tc.body).Return(tc.expectedID, tc.expectedErr).AnyTimes()

		a, _ := mock.Post(context.TODO(), tc.body)

		if a != tc.expectedAuthor {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestPut : test the logic of updating an author
func TestPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := store.NewMockAuthorStorer(ctrl)
	mock := New(mockStore)

	testcases := []struct {
		desc     string
		input    entities.Author
		targetID int

		expected    entities.Author
		expectedErr error
	}{
		{desc: "existing author", input: entities.Author{
			AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/05/1990", PenName: "Dark horse"},
			targetID: 5, expected: entities.Author{AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal",
				DOB: "20/05/1990", PenName: "Dark horse"}, expectedErr: nil,
		},
		{desc: "not existing author", input: entities.Author{
			AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/05/1990", PenName: "Dark horse"},
			targetID: 10, expected: entities.Author{}, expectedErr: errors.New("already exist"),
		},
		{desc: "invalid case", input: entities.Author{
			AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/05/1990", PenName: "Dark horse"},
			targetID: 5, expected: entities.Author{}, expectedErr: errors.New("already exist"),
		},
		{desc: "invalid firstname", input: entities.Author{
			AuthorID: 3, FirstName: "", LastName: "mrinal", DOB: "20/05/1990", PenName: "Dark horse"},
			targetID: 5, expected: entities.Author{}, expectedErr: errors.New("invalid constraints"),
		},
		{desc: "invalid DOB", input: entities.Author{
			AuthorID: 3, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/00/1990", PenName: "Dark horse"},
			targetID: 5, expected: entities.Author{}, expectedErr: errors.New("invalid constraints"),
		},
	}
	author := entities.Author{AuthorID: 5, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/05/1990", PenName: "Dark horse"}

	for _, tc := range testcases {
		if tc.input.AuthorID == 4 && tc.targetID == 10 {
			mockStore.EXPECT().IncludeAuthor(context.TODO(), tc.targetID).Return(author, tc.expectedErr)
		}

		if tc.input.AuthorID == 4 && tc.targetID == 5 {
			mockStore.EXPECT().IncludeAuthor(context.TODO(), tc.targetID).Return(author, nil)
			mockStore.EXPECT().Put(context.TODO(), tc.input, tc.targetID).Return(tc.input.AuthorID, tc.expectedErr)
		}

		author1, _ := mock.Put(context.TODO(), tc.input, tc.targetID)

		if author1 != tc.expected {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestDelete : test logic of deleting an author
func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := store.NewMockAuthorStorer(ctrl)
	mock := New(mockStore)

	testcases := []struct {
		desc     string
		targetID int

		expectedRowsAffected int
		expectedErr          error
	}{
		{"valid authorId", 4, 1, nil},
		{"invalid authorId", -1, 0, errors.New("invalid id")},
		{"error case", 4, 0, errors.New("invalid id")},
	}

	for _, tc := range testcases {
		if tc.targetID == 4 {
			mockStore.EXPECT().Delete(context.TODO(), tc.targetID).Return(tc.expectedRowsAffected, tc.expectedErr)
		}

		id, _ := mock.Delete(context.TODO(), tc.targetID)
		if id != tc.expectedRowsAffected {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}
