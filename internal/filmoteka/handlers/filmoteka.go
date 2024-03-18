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
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

type actor struct {
	srv domain.ActorService
}

func NewActor(srv domain.ActorService) *actor {
	return &actor{srv: srv}
}

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

type film struct {
	srv domain.FilmService
}

func NewFilm(srv domain.FilmService) *film {
	return &film{srv: srv}
}

func (h *film) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.srv.Ping(r.Context())
	if err != nil {
		logger.Logger().Errorln(zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
