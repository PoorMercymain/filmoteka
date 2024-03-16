package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
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
