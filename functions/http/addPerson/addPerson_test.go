package addPerson

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

type Input struct {
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Dob      string `json:"dob"`
	Postcode string `json:"postcode"`
}

type FakeStore struct {
	err bool
}

func (s FakeStore) SaveDoc(ctx context.Context, collection string, data interface{}) (*firestore.DocumentRef, error) {
	if s.err {
		return &firestore.DocumentRef{}, fmt.Errorf("error")
	} else {
		return &firestore.DocumentRef{
			ID: "abc",
		}, nil
	}
}

func TestAddPerson(t *testing.T) {
	tests := []struct {
		name       string
		input      Input
		expected   string
		statusCode int
		store      FakeStore
	}{
		{
			name: "should add a person",
			input: Input{
				Fname:    "Test",
				Lname:    "User",
				Dob:      "1980-01-01T00:00:00.000Z",
				Postcode: "E1 7HQ",
			},
			expected:   "successfully added user. Id: abc",
			statusCode: http.StatusOK,
			store: FakeStore{
				err: false,
			},
		},
		{
			name: "should error if the storage failed",
			input: Input{
				Fname:    "Test",
				Lname:    "User",
				Dob:      "1980-01-01T00:00:00.000Z",
				Postcode: "E1 7HQ",
			},
			expected:   "error saving user: error",
			statusCode: http.StatusInternalServerError,
			store: FakeStore{
				err: true,
			},
		},
		{
			name: "should error if creating the database person failed e.g. incorrect timestamp",
			input: Input{
				Fname:    "Test",
				Lname:    "User",
				Dob:      "abceakfj",
				Postcode: "E1 7HQ",
			},
			expected:   "incorrect time format use: 2006-01-02T15:04:05.000Z",
			statusCode: http.StatusUnprocessableEntity,
			store: FakeStore{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("error marshalling input: %v", err)
			}
			req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			runAddPerson(rr, req, tt.store)

			out, err := ioutil.ReadAll(rr.Result().Body)
			if err != nil {
				t.Fatalf("Error reading response: %v", err)
			}
			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, tt.expected, string(out))
		})
	}
}
