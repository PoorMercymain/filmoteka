package httperrorwriter

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

func WriteError(w http.ResponseWriter, err error, statusCode int, logMsgPrefix string) {
	logger.Logger().Errorln(logMsgPrefix, zap.Error(err))
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		logger.Logger().Errorln(logMsgPrefix, zap.Error(err))
	}
}
