package okta

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestData struct {
	Foo string `json:"foo"`
}

var requestsTests = []struct {
	testServer *httptest.Server
	method     string
	url        string
	body       io.Reader
	target     interface{}
	errCheck   func(t *testing.T, err error)
}{
	// GET request
	{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})),
		"GET",
		"/test",
		nil,
		nil,
		func(t *testing.T, err error) {
			if err != nil {
				t.Error("Expected no error", err)
			}
		},
	},
	// GET request with bad gateway error
	{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		})),
		"GET",
		"/test",
		nil,
		nil,
		func(t *testing.T, err error) {
			if err == nil {
				t.Error("Expected error", err)
			}
		},
	},
	// GET request with not found error
	{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})),
		"GET",
		"/test",
		nil,
		nil,
		func(t *testing.T, err error) {
			if err != ErrNotFound {
				t.Error("Expected error type ErrorNotFound got", err)
			}
		},
	},
	// POST request
	{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, string([]byte(`{"foo":"bar"}`)))
		})),
		"POST",
		"/test",
		bytes.NewBufferString("some bytes"),
		&TestData{},
		func(t *testing.T, err error) {
			if err != nil {
				t.Error("Expected no error got", err)
			}
		},
	},
}

func TestSendRequest(t *testing.T) {
	client := &Client{client: http.DefaultClient}

	for _, tt := range requestsTests {
		err := client.sendRequest(context.Background(),
			tt.method,
			fmt.Sprintf("%s%s", tt.testServer.URL, tt.url),
			tt.body,
			tt.target)
		tt.errCheck(t, err)
		if tt.target != nil {
			testData := tt.target.(*TestData)
			if testData.Foo == "" {
				t.Error("Target is empty")
			}
		}

	}
}
