package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain/mocks"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/service"
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

func testRouter(t *testing.T) *http.ServeMux {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mux := http.NewServeMux()

	ar := mocks.NewMockActorRepository(ctrl)
	as := service.NewActor(ar)
	ah := NewActor(as)

	fr := mocks.NewMockFilmRepository(ctrl)
	fs := service.NewFilm(fr)
	fh := NewFilm(fs)

	aur := mocks.NewMockAuthorizationRepository(ctrl)
	aus := service.NewAuthorization(aur)
	auh := NewAuthorization(aus, "")

	hash, err := bcrypt.GenerateFromPassword([]byte("abc"), bcrypt.DefaultCost)
	require.NoError(t, err)

	ar.EXPECT().CreateActor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, errors.New("")).MaxTimes(1)
	ar.EXPECT().CreateActor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil).MaxTimes(1)
	ar.EXPECT().UpdateActor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrNotFoundInDB).MaxTimes(1)
	ar.EXPECT().UpdateActor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)
	ar.EXPECT().UpdateActor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
	ar.EXPECT().DeleteActor(gomock.Any(), gomock.Any()).Return(appErrors.ErrNotFoundInDB).MaxTimes(1)
	ar.EXPECT().DeleteActor(gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)
	ar.EXPECT().DeleteActor(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
	ar.EXPECT().ReadActors(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("")).MaxTimes(1)
	ar.EXPECT().ReadActors(gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.OutputActor, 0), nil).MaxTimes(1)
	ar.EXPECT().ReadActors(gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.OutputActor, 1), nil).MaxTimes(1)
	fr.EXPECT().CreateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, appErrors.ErrActorNotBornBeforeFilmRelease).MaxTimes(1)
	fr.EXPECT().CreateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, appErrors.ErrActorDoesNotExist).MaxTimes(1)
	fr.EXPECT().CreateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, errors.New("")).MaxTimes(1)
	fr.EXPECT().CreateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil).MaxTimes(1)
	fr.EXPECT().UpdateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrActorNotBornBeforeFilmRelease).MaxTimes(1)
	fr.EXPECT().UpdateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrActorDoesNotExist).MaxTimes(1)
	fr.EXPECT().UpdateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrNotFoundInDB).MaxTimes(1)
	fr.EXPECT().UpdateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)
	fr.EXPECT().UpdateFilm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
	fr.EXPECT().DeleteFilm(gomock.Any(), gomock.Any()).Return(appErrors.ErrNotFoundInDB).MaxTimes(1)
	fr.EXPECT().DeleteFilm(gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)
	fr.EXPECT().DeleteFilm(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
	fr.EXPECT().ReadFilms(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("")).MaxTimes(1)
	fr.EXPECT().ReadFilms(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.OutputFilm, 0), nil).MaxTimes(1)
	fr.EXPECT().ReadFilms(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.OutputFilm, 1), nil).MaxTimes(1)
	fr.EXPECT().FindFilms(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("")).MaxTimes(1)
	fr.EXPECT().FindFilms(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.OutputFilm, 0), nil).MaxTimes(1)
	fr.EXPECT().FindFilms(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]domain.OutputFilm, 1), nil).MaxTimes(1)
	aur.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any()).Return(appErrors.ErrAlreadyRegistered).MaxTimes(1)
	aur.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)
	aur.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	aur.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(false, appErrors.ErrUserNotFound).MaxTimes(1)
	aur.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(false, errors.New("")).MaxTimes(1)
	aur.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(false, nil).MaxTimes(1)
	aur.EXPECT().GetPasswordHash(gomock.Any(), gomock.Any()).Return("", appErrors.ErrUserNotFound).MaxTimes(1)
	aur.EXPECT().GetPasswordHash(gomock.Any(), gomock.Any()).Return("", nil).MaxTimes(1)
	aur.EXPECT().GetPasswordHash(gomock.Any(), gomock.Any()).Return("", errors.New("")).MaxTimes(1)
	aur.EXPECT().GetPasswordHash(gomock.Any(), gomock.Any()).Return(string(hash), nil).MaxTimes(3)
	aur.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(false, appErrors.ErrUserNotFound).MaxTimes(1)
	aur.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(false, errors.New("")).MaxTimes(1)
	aur.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(false, nil).MaxTimes(1)

	mux.Handle("POST /actor", http.HandlerFunc(ah.CreateActor))
	mux.Handle("PUT /actor/{id}", http.HandlerFunc(ah.UpdateActor))
	mux.Handle("DELETE /actor/{id}", http.HandlerFunc(ah.DeleteActor))
	mux.Handle("GET /actors", http.HandlerFunc(ah.ReadActors))

	mux.Handle("POST /film", http.HandlerFunc(fh.CreateFilm))
	mux.Handle("PUT /film/{id}", http.HandlerFunc(fh.UpdateFilm))
	mux.Handle("DELETE /film/{id}", http.HandlerFunc(fh.DeleteFilm))
	mux.Handle("GET /films", http.HandlerFunc(fh.ReadFilms))
	mux.Handle("GET /films/search", http.HandlerFunc(fh.FindFilms))

	mux.Handle("POST /register", http.HandlerFunc(auh.Register))
	mux.Handle("POST /login", http.HandlerFunc(auh.LogIn))

	return mux
}

func request(t *testing.T, ts *httptest.Server, code int, method, content, body, endpoint string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+endpoint, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", content)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, code, resp.StatusCode)

	return resp
}

func TestCreateActor(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"test\", \"gender\":\"male\",\"birthday\":1}",
		},
		{
			"/actor",
			http.MethodPost,
			"",
			http.StatusBadRequest,
			"{\"name\":\"test\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"test\", \"gender\":\"\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"test\", \"gender\":\"male\",\"birthday\":\"\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"abc\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"abc\", \"gender\":\"female\",\"birthday\":\"2020\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"abc\",\"gender\":\"male\",\"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"name\":\"abc\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor",
			http.MethodPost,
			"application/json",
			http.StatusCreated,
			"{\"name\":\"abc\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestUpdateActor(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/actor/",
			http.MethodPut,
			"application/json",
			http.StatusNotFound,
			"{\"name\":\"test\", \"gender\":\"male\",\"birthday\":\"\"}",
		},
		{
			"/actor/abc",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"\", \"gender\":\"male\",\"birthday\":\"2020-01-02\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"\", \"gender\":\"\",\"birthday\":\"\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"\", \"gender\":\"abc\",\"birthday\":\"\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"\", \"gender\":\"\",\"birthday\":\"abc\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"name\":\"abc\",\"name\":\"abc\",\"gender\":\"\",\"birthday\":\"\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusNotFound,
			"{\"name\":\"abc\", \"gender\":\"\",\"birthday\":\"\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusInternalServerError,
			"{\"name\":\"abc\", \"gender\":\"\",\"birthday\":\"\"}",
		},
		{
			"/actor/1",
			http.MethodPut,
			"application/json",
			http.StatusNoContent,
			"{\"name\":\"abc\", \"gender\":\"\",\"birthday\":\"\"}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestDeleteActor(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/actor/",
			http.MethodDelete,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/actor/abc",
			http.MethodDelete,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/actor/1",
			http.MethodDelete,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/actor/1",
			http.MethodDelete,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/actor/1",
			http.MethodDelete,
			"",
			http.StatusNoContent,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestReadActors(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/actors?page=a",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/actors?limit=a",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/actors?page=0",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/actors?limit=0",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/actors?limit=101",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/actors?page=1&limit=10",
			http.MethodGet,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/actors?page=1&limit=10",
			http.MethodGet,
			"",
			http.StatusNoContent,
			"",
		},
		{
			"/actors?page=1&limit=10",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestCreateFilm(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/film",
			http.MethodPost,
			"",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"abc\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": \"\",\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [\"3\", \"4\"]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"a\",\"releaseDate\": \"2021-04-13\",\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"a\",\"releaseDate\": \"2021-04-13\",\"rating\":12,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2019-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusNotFound,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2019-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2019-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film",
			http.MethodPost,
			"application/json",
			http.StatusCreated,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2019-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestUpdateFilm(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/film/",
			http.MethodPut,
			"",
			http.StatusNotFound,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/abc",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"",
			http.StatusBadRequest,
			"{\"title\":\"abc\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\"releaseDate\": \"2021-04-13\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2\",\"rating\": 5.7,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 12,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusBadRequest,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 1,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusNotFound,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 1,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusNotFound,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 1,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusInternalServerError,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 1,\"actorIDs\": [3, 4]}",
		},
		{
			"/film/1",
			http.MethodPut,
			"application/json",
			http.StatusNoContent,
			"{\"title\":\"a\",\"description\": \"test\",\"releaseDate\": \"2021-04-13\",\"rating\": 1,\"actorIDs\": [3, 4]}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestDeleteFilm(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/film/a",
			http.MethodDelete,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/film/1",
			http.MethodDelete,
			"",
			http.StatusNotFound,
			"",
		},
		{
			"/film/1",
			http.MethodDelete,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/film/1",
			http.MethodDelete,
			"",
			http.StatusNoContent,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestReadFilms(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/films?field=abc",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films?order=abc",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films?page=abc",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films?limit=abc",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films?page=0",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films?limit=0",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films?limit=101",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films",
			http.MethodGet,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/films",
			http.MethodGet,
			"",
			http.StatusNoContent,
			"",
		},
		{
			"/films",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestFindFilms(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/films/search?page=abc",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films/search?limit=abc",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films/search?page=0",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films/search?limit=101",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films/search",
			http.MethodGet,
			"",
			http.StatusBadRequest,
			"",
		},
		{
			"/films/search?title=a",
			http.MethodGet,
			"",
			http.StatusInternalServerError,
			"",
		},
		{
			"/films/search?title=a",
			http.MethodGet,
			"",
			http.StatusNoContent,
			"",
		},
		{
			"/films/search?title=a",
			http.MethodGet,
			"",
			http.StatusOK,
			"",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestRegister(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/register",
			http.MethodPost,
			"",
			http.StatusBadRequest,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"login\":\"abc\",\"password\":\"abc\"",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"login\":1,\"password\":\"abc\"}",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusConflict,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/register",
			http.MethodPost,
			"application/json",
			http.StatusOK,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
	}

	for _, testCase := range testTable {
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}

func TestLogIn(t *testing.T) {
	ts := httptest.NewServer(testRouter(t))

	defer ts.Close()

	var testTable = []struct {
		endpoint string
		method   string
		content  string
		code     int
		body     string
	}{
		{
			"/login",
			http.MethodPost,
			"",
			http.StatusBadRequest,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"login\":\"abc\",\"password\":\"abc\"",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusBadRequest,
			"{\"login\":1,\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusUnauthorized,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusUnauthorized,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusInternalServerError,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
		{
			"/login",
			http.MethodPost,
			"application/json",
			http.StatusOK,
			"{\"login\":\"abc\",\"password\":\"abc\"}",
		},
	}

	for i, testCase := range testTable {
		logger.Logger().Infoln(i)
		resp := request(t, ts, testCase.code, testCase.method, testCase.content, testCase.body, testCase.endpoint)
		resp.Body.Close()
	}
}
