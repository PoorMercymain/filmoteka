package jsonhttpvalidator

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testResponseWriter struct {
	http.ResponseWriter
}

func (trw testResponseWriter) WriteHeader(statusCode int) {}

func (trw testResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}

type testReqBody struct {
	io.ReadCloser
}

func (trb *testReqBody) Read(p []byte) (int, error) {
	n, err := trb.ReadCloser.Read(p)
	if err != nil {
		err = errors.New("")
	}

	return n, err
}

func TestValidate(t *testing.T) {
	trw := testResponseWriter{}

	r, err := http.NewRequest("POST", "/", strings.NewReader(""))
	require.NoError(t, err)

	r.Header.Set("Content-Type", "test")
	err = ValidateJSONRequest(trw, r, "")
	require.Error(t, err)

	r, err = http.NewRequest("POST", "/", strings.NewReader("{\"id\":0,\"id\":0}"))
	require.NoError(t, err)

	r.Header.Set("Content-Type", "application/json")
	err = ValidateJSONRequest(trw, r, "")
	require.Error(t, err)

	r, err = http.NewRequest("POST", "/", strings.NewReader("{\"id\":0}"))
	require.NoError(t, err)

	r.Header.Set("Content-Type", "application/json")
	err = ValidateJSONRequest(trw, r, "")
	require.NoError(t, err)

	r, err = http.NewRequest("POST", "/", strings.NewReader("{\"id\":0}"))
	require.NoError(t, err)

	r.Header.Set("Content-Type", "application/json")
	r.Body = &testReqBody{r.Body}
	err = ValidateJSONRequest(trw, r, "")
	require.Error(t, err)
}
