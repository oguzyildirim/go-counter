package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Oguzyildirim/go-counter/internal/counter/rest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var Create func() error

var Find func() (string, error)

var Handle func(next http.Handler) http.Handler

type counterServiceMock struct{}

type middlewareServiceMock struct{}

func (svc counterServiceMock) Create() error {
	return Create()
}

func (svc counterServiceMock) Find() (string, error) {
	return Find()
}

func (m middlewareServiceMock) Handle(next http.Handler) http.Handler {
	return Handle(next)
}

func TestCreate(t *testing.T) {
	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		input  []byte
		output output
	}{
		{
			"OK: 201",
			nil,
			output{
				http.StatusCreated,
				&struct{}{},
				&struct{}{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mockSvc := counterServiceMock{}
			mockMiddleware := middlewareServiceMock{}
			Create = func() error {
				return nil
			}
			Find = func() (string, error) {
				return "", nil
			}
			Handle = func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			}

			router := http.NewServeMux()

			rest.NewCounterHandler(mockSvc, mockMiddleware).Register(router)

			res := doRequest(router, httptest.NewRequest(http.MethodGet, "/", nil))

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		input  []byte
		output output
	}{
		{
			"OK: 201",
			nil,
			output{
				http.StatusOK,
				&rest.FindCounterResponse{Count: "10"},
				&rest.FindCounterResponse{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mockSvc := counterServiceMock{}
			mockMiddleware := middlewareServiceMock{}
			Create = func() error {
				return nil
			}
			Find = func() (string, error) {
				return "10", nil
			}
			Handle = func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			}

			router := http.NewServeMux()

			rest.NewCounterHandler(mockSvc, mockMiddleware).Register(router)

			res := doRequest(router, httptest.NewRequest(http.MethodGet, "/count", nil))

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

func doRequest(router *http.ServeMux, req *http.Request) *http.Response {
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr.Result()
}

type test struct {
	expected interface{}
	target   interface{}
}

func assertResponse(t *testing.T, res *http.Response, test test) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(test.target); err != nil {
		t.Fatalf("couldn't decode %s", err)
	}
	defer res.Body.Close()

	// https://github.com/stretchr/testify/issues/535
	// go-cmp requires that for each struct in a comparison that has any unexported fields you must explicitly use either AllowUnexporterted or IgnoreUnexported.
	//This requirement seems like one of the advantages of using go-cmp instead of reflect.DeepEqual.
	/* 	if !reflect.DeepEqual(test.expected, test.target) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(test.expected, test.target, cmpopts.IgnoreUnexported()))
	} */
	if !cmp.Equal(test.expected, test.target, cmpopts.IgnoreUnexported()) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(test.expected, test.target, cmpopts.IgnoreUnexported()))
	}
}
