package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
	httperrorwriter "github.com/PoorMercymain/filmoteka/pkg/http-error-writer"
	jsonhttpvalidator "github.com/PoorMercymain/filmoteka/pkg/json-http-validator"
	"github.com/PoorMercymain/filmoteka/pkg/jwt"
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

type actor struct {
	srv domain.ActorService
}

func NewActor(srv domain.ActorService) *actor {
	return &actor{srv: srv}
}

// @Tags Actors
// @Summary Запрос добавления актера в БД
// @Description Запрос для добавления информации об актере в БД
// @Accept json
// @Produce json
// @Param input body domain.Actor true "информация об актере"
// @Success 201
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /actor [post]
func (h *actor) CreateActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.CreateActor():"

	err := jsonhttpvalidator.ValidateJSONRequest(w, r, logErrPrefix)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var actor domain.Actor
	if err = d.Decode(&actor); err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	if actor.Name == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoNameProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	var gender bool
	if actor.Gender == "male" {
		gender = domain.Male
	} else if actor.Gender == "female" {
		gender = domain.Female
	} else {
		httperrorwriter.WriteError(w, appErrors.ErrUnknownGender, http.StatusBadRequest, logErrPrefix)
		return
	}

	birthday, err := time.Parse(time.DateOnly, actor.Birthday)
	if err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	id, err := h.srv.CreateActor(r.Context(), actor.Name, gender, birthday)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	e := json.NewEncoder(w)
	err = e.Encode(domain.ID{ID: id})
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}

// @Tags Actors
// @Summary Запрос обновления актера в БД
// @Description Запрос для обновления информации об актере в БД, как полностью, так и частичного
// @Accept json
// @Param input body domain.Actor true "информация об актере"
// @Param id path int true "id актера" Example(1)
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /actor/{id} [put]
func (h *actor) UpdateActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.UpdateActor():"

	idStr := r.PathValue("id")

	if idStr == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoIDProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrIDIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	err = jsonhttpvalidator.ValidateJSONRequest(w, r, logErrPrefix)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var actor domain.Actor
	if err = d.Decode(&actor); err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	if actor.Name == "" && actor.Gender == "" && actor.Birthday == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNothingProvidedInJSON, http.StatusBadRequest, logErrPrefix)
		return
	}

	var gender bool
	var genderPtr *bool
	if actor.Gender != "" {
		if actor.Gender == "male" {
			gender = domain.Male
			genderPtr = &gender
		} else if actor.Gender == "female" {
			gender = domain.Female
			genderPtr = &gender
		} else {
			httperrorwriter.WriteError(w, appErrors.ErrUnknownGender, http.StatusBadRequest, logErrPrefix)
			return
		}
	}

	var birthday time.Time
	if actor.Birthday != "" {
		birthday, err = time.Parse(time.DateOnly, actor.Birthday)
		if err != nil {
			httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
			return
		}
	}

	err = h.srv.UpdateActor(r.Context(), id, actor.Name, genderPtr, birthday)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFoundInDB) {
			httperrorwriter.WriteError(w, appErrors.ErrNotFoundInDB, http.StatusNotFound, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Actors
// @Summary Запрос удаления актера из БД
// @Description Запрос для удаления информации об актере из БД
// @Param id path int true "id актера" Example(1)
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /actor/{id} [delete]
func (h *actor) DeleteActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.DeleteActor():"

	idStr := r.PathValue("id")

	if idStr == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoIDProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrIDIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	err = h.srv.DeleteActor(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFoundInDB) {
			httperrorwriter.WriteError(w, appErrors.ErrNotFoundInDB, http.StatusNotFound, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Actors
// @Summary Запрос получения списка актеров из БД
// @Description Запрос для получения списка актеров из БД, для каждого актера также выводится список фильмов с его участием, предусмотрена пагинация
// @Produce json
// @Param page query int false "номер страницы, начинается с 1 (по умолчанию 1)" Example(1)
// @Param limit query int false "максимальное число актеров на странице, в диапазоне [1, 100] (по умолчанию 15)" Example(1)
// @Success 200
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /actors [get]
func (h *actor) ReadActors(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.ReadActors():"

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	if pageStr == "" {
		pageStr = "1"
	}

	if limitStr == "" {
		limitStr = "15"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrPageInNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrLimitIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	if page < 1 {
		httperrorwriter.WriteError(w, appErrors.ErrPageNumberIsTooSmall, http.StatusBadRequest, logErrPrefix)
		return
	}

	if limit < 1 || limit > 100 {
		httperrorwriter.WriteError(w, appErrors.ErrLimitParameterNotInCorrectRange, http.StatusBadRequest, logErrPrefix)
		return
	}

	actors, err := h.srv.ReadActors(r.Context(), page, limit)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	if len(actors) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	err = e.Encode(actors)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}

type film struct {
	srv domain.FilmService
}

func NewFilm(srv domain.FilmService) *film {
	return &film{srv: srv}
}

// @Tags Films
// @Summary Запрос добавления информации о фильме в БД
// @Description Запрос для добавления информации о фильме в БД
// @Accept json
// @Produce json
// @Param input body domain.Film true "информация о фильме"
// @Success 201
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /film [post]
func (h *film) CreateFilm(w http.ResponseWriter, r *http.Request) {
	const (
		titleLimit       = 150
		descriptionLimit = 1000
		minRating        = 0
		maxRating        = 10
	)

	defer r.Body.Close()
	const logErrPrefix = "handlers.CreateFilm():"

	err := jsonhttpvalidator.ValidateJSONRequest(w, r, logErrPrefix)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var film domain.Film
	if err = d.Decode(&film); err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	if film.Title == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoTitleProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	if len([]rune(film.Title)) > titleLimit {
		httperrorwriter.WriteError(w, appErrors.ErrTitleTooLong, http.StatusBadRequest, logErrPrefix)
		return
	}

	if film.Description == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoDescriptionProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	if len([]rune(film.Description)) > descriptionLimit {
		httperrorwriter.WriteError(w, appErrors.ErrTitleTooLong, http.StatusBadRequest, logErrPrefix)
		return
	}

	releaseDate, err := time.Parse(time.DateOnly, film.ReleaseDate)
	if err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	if film.Rating == nil {
		httperrorwriter.WriteError(w, appErrors.ErrNoRatingValue, http.StatusBadRequest, logErrPrefix)
		return
	}

	if *film.Rating < minRating || *film.Rating > maxRating {
		httperrorwriter.WriteError(w, appErrors.ErrWrongRatingValue, http.StatusBadRequest, logErrPrefix)
		return
	}

	if film.Actors == nil {
		film.Actors = make([]int, 0)
	}

	id, err := h.srv.CreateFilm(r.Context(), film.Title, film.Description, releaseDate, *film.Rating, film.Actors)
	if err != nil {
		if errors.Is(err, appErrors.ErrActorNotBornBeforeFilmRelease) {
			httperrorwriter.WriteError(w, appErrors.ErrActorNotBornBeforeFilmRelease, http.StatusBadRequest, logErrPrefix)
			return
		}

		if errors.Is(err, appErrors.ErrActorDoesNotExist) {
			httperrorwriter.WriteError(w, appErrors.ErrActorDoesNotExist, http.StatusNotFound, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	e := json.NewEncoder(w)
	err = e.Encode(domain.ID{ID: id})
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}

// @Tags Films
// @Summary Запрос обновления информации о фильме
// @Description Запрос для обновления информации о фильме, как полного, так и частичного
// @Accept json
// @Param input body domain.Film true "информация о фильме, если не убрать из запроса поле actorIDs, его значение заменит актеров фильма в БД"
// @Param id path int true "id фильма" Example(1)
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /film/{id} [put]
func (h *film) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	const (
		titleLimit       = 150
		descriptionLimit = 1000
		minRating        = 0
		maxRating        = 10
	)

	defer r.Body.Close()
	const logErrPrefix = "handlers.UpdateFilm():"

	idStr := r.PathValue("id")

	if idStr == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoIDProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrIDIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	err = jsonhttpvalidator.ValidateJSONRequest(w, r, logErrPrefix)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var film domain.Film
	if err = d.Decode(&film); err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	if film.Actors == nil && film.Description == "" && film.Rating == nil && film.ReleaseDate == "" && film.Title == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNothingProvidedInJSON, http.StatusBadRequest, logErrPrefix)
		return
	}

	if len([]rune(film.Title)) > titleLimit {
		httperrorwriter.WriteError(w, appErrors.ErrTitleTooLong, http.StatusBadRequest, logErrPrefix)
		return
	}

	if len([]rune(film.Description)) > descriptionLimit {
		httperrorwriter.WriteError(w, appErrors.ErrTitleTooLong, http.StatusBadRequest, logErrPrefix)
		return
	}

	var releaseDate time.Time
	if film.ReleaseDate != "" {
		releaseDate, err = time.Parse(time.DateOnly, film.ReleaseDate)
		if err != nil {
			httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
			return
		}
	}

	if film.Rating != nil {
		if *film.Rating < minRating || *film.Rating > maxRating {
			httperrorwriter.WriteError(w, appErrors.ErrWrongRatingValue, http.StatusBadRequest, logErrPrefix)
			return
		}
	}

	err = h.srv.UpdateFilm(r.Context(), id, film.Title, film.Description, releaseDate, film.Rating, film.Actors)
	if err != nil {
		if errors.Is(err, appErrors.ErrActorNotBornBeforeFilmRelease) {
			httperrorwriter.WriteError(w, appErrors.ErrActorNotBornBeforeFilmRelease, http.StatusBadRequest, logErrPrefix)
			return
		}

		if errors.Is(err, appErrors.ErrActorDoesNotExist) {
			httperrorwriter.WriteError(w, appErrors.ErrActorDoesNotExist, http.StatusNotFound, logErrPrefix)
			return
		}

		if errors.Is(err, appErrors.ErrNotFoundInDB) {
			httperrorwriter.WriteError(w, appErrors.ErrNotFoundInDB, http.StatusNotFound, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Films
// @Summary Запрос удаления фильма из БД
// @Description Запрос для удаления информации о фильме из БД
// @Param id path int true "id фильма" Example(1)
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /film/{id} [delete]
func (h *film) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.DeleteFilm():"

	idStr := r.PathValue("id")

	if idStr == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoIDProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrIDIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	err = h.srv.DeleteFilm(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFoundInDB) {
			httperrorwriter.WriteError(w, appErrors.ErrNotFoundInDB, http.StatusNotFound, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Films
// @Summary Запрос получения списка фильмов из БД
// @Description Запрос для получения списка фильмов из БД, для каждого фильма также выводится список фильмов с его участием, предусмотрена пагинация, по умолчанию сортируется по убыванию рейтинга
// @Produce json
// @Param field query string false "поле для сортировки (release_date, rating, title, по умолчанию - rating)" Example(title)
// @Param order query string false "поле для порядка сортировки (desc - по убыванию, asc - по возрастанию, по умолчанию - desc)" Example(desc)
// @Param page query int false "номер страницы, начинается с 1 (по умолчанию 1)" Example(1)
// @Param limit query int false "максимальное число актеров на странице, в диапазоне [1, 100] (по умолчанию 15)" Example(1)
// @Success 200
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /films [get]
func (h *film) ReadFilms(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.ReadFilms():"

	field := r.URL.Query().Get("field")
	order := r.URL.Query().Get("order")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if field == "" {
		field = "rating"
	}

	if order == "" {
		order = "desc"
	}

	if pageStr == "" {
		pageStr = "1"
	}

	if limitStr == "" {
		limitStr = "15"
	}

	if field != "title" && field != "rating" && field != "release_date" {
		httperrorwriter.WriteError(w, appErrors.ErrUnknownSortField, http.StatusBadRequest, logErrPrefix)
		return
	}

	if order != "desc" && order != "asc" {
		httperrorwriter.WriteError(w, appErrors.ErrUnknownOrder, http.StatusBadRequest, logErrPrefix)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrPageInNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrLimitIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	if page < 1 {
		httperrorwriter.WriteError(w, appErrors.ErrPageNumberIsTooSmall, http.StatusBadRequest, logErrPrefix)
		return
	}

	if limit < 1 || limit > 100 {
		httperrorwriter.WriteError(w, appErrors.ErrLimitParameterNotInCorrectRange, http.StatusBadRequest, logErrPrefix)
		return
	}

	films, err := h.srv.ReadFilms(r.Context(), field, order, page, limit)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	if len(films) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	err = e.Encode(films)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}

// @Tags Films
// @Summary Запрос поиска фильмов в БД
// @Description Запрос для поиска фильмов в БД по фрагменту названия фильма и/или имени актера, по умолчанию выдает 1 самый подходящий фильм, для успешного запроса надо указать хотя бы один из фрагментов
// @Produce json
// @Param title query string false "фрагмент названия фильма для поиска" Example(film)
// @Param order query string false "фрагмент имени актера для поиска" Example(Val)
// @Param page query int false "номер страницы, начинается с 1 (по умолчанию 1)" Example(1)
// @Param limit query int false "максимальное число актеров на странице, в диапазоне [1, 100] (по умолчанию 1)" Example(1)
// @Success 200
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /films/search [get]
func (h *film) FindFilms(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.FindFilms():"

	filmTitleFragment := r.URL.Query().Get("title")
	actorNameFragment := r.URL.Query().Get("name")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	if pageStr == "" {
		pageStr = "1"
	}

	if limitStr == "" {
		limitStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrPageInNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		httperrorwriter.WriteError(w, appErrors.ErrLimitIsNotANumber, http.StatusBadRequest, logErrPrefix)
		return
	}

	if page < 1 {
		httperrorwriter.WriteError(w, appErrors.ErrPageNumberIsTooSmall, http.StatusBadRequest, logErrPrefix)
		return
	}

	if limit < 1 || limit > 100 {
		httperrorwriter.WriteError(w, appErrors.ErrLimitParameterNotInCorrectRange, http.StatusBadRequest, logErrPrefix)
		return
	}

	if filmTitleFragment == "" && actorNameFragment == "" {
		httperrorwriter.WriteError(w, appErrors.ErrNoFragmentsProvided, http.StatusBadRequest, logErrPrefix)
		return
	}

	films, err := h.srv.FindFilms(r.Context(), filmTitleFragment, actorNameFragment, page, limit)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	if len(films) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	err = e.Encode(films)
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}

type authorization struct {
	srv    domain.AuthorizationService
	JWTKey string
}

func NewAuthorization(srv domain.AuthorizationService, jwtKey string) *authorization {
	return &authorization{srv: srv, JWTKey: jwtKey}
}

// @Tags Auth
// @Summary Запрос регистрации в filmoteka
// @Description Запрос для регистрации в сервисе, производится регистрация обычного пользователя (если нужен админ, надо задать соответствующее поле в БД в таблице auth и заново получить токен через login) и выдается JWT (можно указать в заголовке Authorization) на 24 часа (также записывается в Cookie)
// @Accept json
// @Produce json
// @Param input body domain.AuthorizationData true "аутентификационные данные"
// @Success 201
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /register [post]
func (h *authorization) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.Register():"

	err := jsonhttpvalidator.ValidateJSONRequest(w, r, logErrPrefix)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var authData domain.AuthorizationData
	if err = d.Decode(&authData); err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	err = h.srv.Register(r.Context(), authData.Login, authData.Password)
	if err != nil {
		if errors.Is(err, appErrors.ErrAlreadyRegistered) {
			httperrorwriter.WriteError(w, appErrors.ErrAlreadyRegistered, http.StatusConflict, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	isAdmin, err := h.srv.IsAdmin(r.Context(), authData.Login)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			httperrorwriter.WriteError(w, appErrors.ErrUserNotFound, http.StatusInternalServerError, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	tokenStr, err := jwt.CreateJWT(isAdmin, []byte(h.JWTKey), time.Now().Add(24*time.Hour))
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	cookie := http.Cookie{
		Name:   "authToken",
		Value:  tokenStr,
		MaxAge: 86400,
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	e := json.NewEncoder(w)
	err = e.Encode(domain.Token{Token: tokenStr})
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}

// @Tags Auth
// @Summary Запрос получения токена для авторизации
// @Description Запрос для получения JWT в Cookie и теле ответа
// @Accept json
// @Produce json
// @Param input body domain.AuthorizationData true "аутентификационные данные"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /login [post]
func (h *authorization) LogIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const logErrPrefix = "handlers.LogIn():"

	err := jsonhttpvalidator.ValidateJSONRequest(w, r, logErrPrefix)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var authData domain.AuthorizationData
	if err = d.Decode(&authData); err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return
	}

	err = h.srv.CheckAuth(r.Context(), authData.Login, authData.Password)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			httperrorwriter.WriteError(w, appErrors.ErrUserNotFound, http.StatusUnauthorized, logErrPrefix)
			return
		}

		if errors.Is(err, appErrors.ErrWrongPassword) {
			httperrorwriter.WriteError(w, appErrors.ErrWrongPassword, http.StatusUnauthorized, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	isAdmin, err := h.srv.IsAdmin(r.Context(), authData.Login)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			httperrorwriter.WriteError(w, appErrors.ErrUserNotFound, http.StatusInternalServerError, logErrPrefix)
			return
		}

		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	logger.Logger().Infoln(isAdmin)
	tokenStr, err := jwt.CreateJWT(isAdmin, []byte(h.JWTKey), time.Now().Add(24*time.Hour))
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
		httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
		return
	}

	cookie := http.Cookie{
		Name:   "authToken",
		Value:  tokenStr,
		MaxAge: 86400,
	}

	http.SetCookie(w, &cookie)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	err = e.Encode(domain.Token{Token: tokenStr})
	if err != nil {
		logger.Logger().Errorln(logErrPrefix, zap.Error(err))
	}
}
