package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
	jsonduplicatechecker "github.com/PoorMercymain/filmoteka/pkg/json-duplicate-checker"
	jsonmimechecker "github.com/PoorMercymain/filmoteka/pkg/json-mime-checker"
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

type filmoteka struct {
	srv domain.FilmotekaService
}

func New(srv domain.FilmotekaService) *filmoteka {
	return &filmoteka{srv: srv}
}

func (h *filmoteka) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.srv.Ping(r.Context())
	if err != nil {
		logger.Logger().Errorln(zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *filmoteka) CreateActor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !jsonmimechecker.IsJSONContentTypeCorrect(r) {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(appErrors.ErrWrongMIME))
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(appErrors.ErrWrongMIME.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	bytesToCheck, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	reader := bytes.NewReader(bytes.Clone(bytesToCheck))

	err = jsonduplicatechecker.CheckDuplicatesInJSON(json.NewDecoder(reader), nil)
	if err != nil {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bytesToCheck))

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var actor domain.Actor
	if err := d.Decode(&actor); err != nil {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	var gender bool
	if actor.Gender == "male" {
		gender = domain.Male
	} else if actor.Gender == "female" {
		gender = domain.Female
	} else {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(appErrors.ErrUnknownGender))
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(appErrors.ErrUnknownGender.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	birthday, err := time.Parse(time.DateOnly, actor.Birthday)
	if err != nil {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	id, err := h.srv.CreateActor(r.Context(), actor.Name, gender, birthday)
	if err != nil {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	err = e.Encode(domain.ID{ID: id})
	if err != nil {
		logger.Logger().Errorln("handlers.CreateActor():", zap.Error(err))
	}
}
