package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain/mocks"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/handlers"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/service"
	"github.com/PoorMercymain/filmoteka/pkg/jwt"
)

func testRouter(t *testing.T) *http.ServeMux {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mux := http.NewServeMux()

	aur := mocks.NewMockAuthorizationRepository(ctrl)
	aus := service.NewAuthorization(aur)
	auh := handlers.NewAuthorization(aus, "")

	mux.Handle("GET /admin", AdminRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), auh.JWTKey))
	mux.Handle("GET /user", AuthorizationRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), auh.JWTKey))

	return mux
}

func request(t *testing.T, ts *httptest.Server, code int, method, content, body, endpoint, authorization, cookie string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+endpoint, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", content)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: "authToken", Value: cookie})
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, code, resp.StatusCode)

	return resp
}

func TestAdminRequired(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	tokenStrNoAdmin, err := jwt.CreateJWT(false, []byte(""), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	tokenStrAdmin, err := jwt.CreateJWT(true, []byte(""), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	wrongToken, err := jwt.CreateJWT(true, []byte("abcd"), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	var testTable = []struct {
		endpoint      string
		method        string
		content       string
		code          int
		body          string
		authorization string
		cookie string
	}{
		{
			"/admin",
			http.MethodGet,
			"",
			http.StatusUnauthorized,
			"",
			"",
			"",
		},
		{
			"/admin",
			http.MethodGet,
			"",
			http.StatusForbidden,
			"",
			tokenStrNoAdmin,
			"",
		},
		{
			"/admin",
			http.MethodGet,
			"",
			http.StatusUnauthorized,
			"",
			wrongToken,
			"",
		},
		{
			"/admin",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
			tokenStrAdmin,
			"",
		},
		{
			"/admin",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
			"",
			tokenStrAdmin,
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint, testCase.authorization, testCase.cookie)
		resp.Body.Close()
	}
}

func TestAuthorizationRequired(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	tokenStrNoAdmin, err := jwt.CreateJWT(false, []byte(""), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	tokenStrAdmin, err := jwt.CreateJWT(true, []byte(""), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	wrongToken, err := jwt.CreateJWT(true, []byte("abcd"), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	var testTable = []struct {
		endpoint      string
		method        string
		content       string
		code          int
		body          string
		authorization string
		cookie string
	}{
		{
			"/user",
			http.MethodGet,
			"",
			http.StatusUnauthorized,
			"",
			"",
			"",
		},
		{
			"/user",
			http.MethodGet,
			"",
			http.StatusUnauthorized,
			"",
			wrongToken,
			"",
		},
		{
			"/user",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
			tokenStrNoAdmin,
			"",
		},
		{
			"/user",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
			tokenStrAdmin,
			"",
		},
		{
			"/user",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
			"",
			tokenStrAdmin,
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint, testCase.authorization, testCase.cookie)
		resp.Body.Close()
	}
}
