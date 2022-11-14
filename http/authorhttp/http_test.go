package authorhttp

import (
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"testing"

	"net/http"
	"net/http/httptest"

	"projects/GoLang-Interns-2022/authorbook/entities"
	"projects/GoLang-Interns-2022/authorbook/service"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

// TestPost : to test Post handler
func TestPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockAuthorService(ctrl)
	mock := New(mockService)

	testcases := []struct {
		desc  string
		input entities.Author

		expected       interface{}
		expectedStatus int
		expectedErr    error
	}{
		{desc: "valid case:", input: entities.Author{
			AuthorID: 3, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/01/1990", PenName: "Dark horse"},
			expected: entities.Author{
				AuthorID: 3, FirstName: "nilotpal", LastName: "mrinal", DOB: "20/01/1990", PenName: "Dark horse"},
			expectedStatus: http.StatusCreated, expectedErr: nil,
		},
		//{desc: "returning error from svc", input: entities.Author{AuthorID: 4, FirstName: "nilotpal", LastName: "mrinal",
		//	DOB: "20/01/1990", PenName: "Dark horse"}, expected: nil,
		//	expectedStatus: http.StatusBadRequest, expectedErr: errors.New("not valid constraints"),
		//},
		{desc: "unmarshalling error ", input: entities.Author{}, expected: nil,
			expectedStatus: http.StatusBadRequest, expectedErr: errors.New("invalid character 'h' looking for beginning of value"),
		},
	}

	for _, tc := range testcases {
		data, _ := json.Marshal(tc.input)

		if tc.desc == "unmarshalling error " {
			data = []byte("hello")
		}

		k := gofr.New()
		r := httptest.NewRequest("POST", "localhost:8000/author", bytes.NewReader(data))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, k)

		if tc.input.AuthorID == 4 {
			mockService.EXPECT().Post(gomock.Any(), tc.input).Return(tc.expected, tc.expectedErr)
		} else if tc.input.AuthorID == 3 {
			mockService.EXPECT().Post(gomock.Any(), tc.input).Return(tc.expected, tc.expectedErr)
		}

		result, _ := mock.Post(ctx)

		if tc.expected != result {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestPut : to test the put handler
func TestPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockAuthorService(ctrl)
	mock := New(mockService)

	testcases := []struct {
		desc           string
		input          entities.Author
		TargetID       string
		expected       entities.Author
		expectedStatus int
		expectedErr    error
	}{
		{desc: "valid case:", input: entities.Author{
			AuthorID: 3, FirstName: "amit", LastName: "kumar", DOB: "20/01/1990", PenName: "Dark horse"},
			TargetID: "4", expected: entities.Author{AuthorID: 3, FirstName: "amit",
				LastName: "kumar", DOB: "20/01/1990", PenName: "Dark horse"}, expectedStatus: http.StatusCreated,
			expectedErr: nil,
		},
		{desc: "strconv error", input: entities.Author{AuthorID: 3, FirstName: "kumar", LastName: "vis",
			DOB: "20/01/1990", PenName: "Dark horse"}, expected: entities.Author{},
			expectedStatus: http.StatusBadRequest, expectedErr: nil,
		},
		{desc: "unmarshalling error ", input: entities.Author{}, expected: entities.Author{},
			expectedStatus: http.StatusBadRequest, expectedErr: nil,
		},
		{desc: "error from svc layer", input: entities.Author{}, TargetID: "5", expected: entities.Author{},
			expectedStatus: http.StatusNotFound, expectedErr: errors.New("invalid error"),
		},
	}

	k := gofr.New()
	for _, tc := range testcases {
		data, err := json.Marshal(tc.input)
		if err != nil {
			log.Print(err)
		}

		if tc.desc == "unmarshalling error " {
			data = []byte("hello")
		}

		r := httptest.NewRequest("PUT", "localhost:8000/author/{id}"+tc.TargetID, bytes.NewReader(data))
		r = mux.SetURLVars(r, map[string]string{"id": tc.TargetID})
		w := httptest.NewRecorder()
		id, _ := strconv.Atoi(tc.TargetID)

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, k)

		mockService.EXPECT().Put(ctx, tc.input, id).Return(tc.expected, tc.expectedErr).AnyTimes()

		_, err = mock.Put(ctx)

		//res := w.Result()
		if tc.expectedErr != err {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

// TestDelete : to test the delete handler
func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := service.NewMockAuthorService(ctrl)
	mock := New(mockService)

	testcases := []struct {
		desc   string
		target string

		expectedStatus int
		expectedErr    error
	}{
		{"valid authorId", "4", http.StatusNoContent, nil},
		{"invalid authorId", "-3", http.StatusBadRequest, errors.New("invalid")},
		{desc: "invalid authorId", expectedStatus: http.StatusBadRequest, expectedErr: errors.New("invalid")},
	}

	k := gofr.New()
	for _, tc := range testcases {
		r := httptest.NewRequest("DELETE", "localhost:8000/author/{id}"+tc.target, nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.target})
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, k)
		id, err := strconv.Atoi(tc.target)
		if err != nil {
			log.Print(err)
		}

		mockService.EXPECT().Delete(ctx, id).Return(tc.expectedStatus, tc.expectedErr).AnyTimes()

		_, err = mock.Delete(ctx)
		if tc.expectedErr != err {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}
