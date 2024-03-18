package middleware

import (
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	httperrorwriter "github.com/PoorMercymain/filmoteka/pkg/http-error-writer"
	"github.com/PoorMercymain/filmoteka/pkg/jwt"
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

func AdminRequired(next http.Handler, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const logErrPrefix = "middleware.AdminRequired():"

		authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if authToken == "" {
			cookie, err := r.Cookie("authToken")
			if err != nil {
				httperrorwriter.WriteError(w, appErrors.ErrNoTokenProvided, http.StatusUnauthorized, logErrPrefix)
				return
			}

			authToken = cookie.Value
		}

		isAdmin, err := jwt.CheckIsAdminInJWT(authToken, jwtKey)
		if err != nil {
			if errors.Is(err, appErrors.ErrTokenIsInvalid) {
				httperrorwriter.WriteError(w, appErrors.ErrTokenIsInvalid, http.StatusUnauthorized, logErrPrefix)
				return
			}

			logger.Logger().Errorln(logErrPrefix, zap.Error(err))
			httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
			return
		}

		if !isAdmin {
			httperrorwriter.WriteError(w, appErrors.ErrAdminRequired, http.StatusForbidden, logErrPrefix)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthorizationRequired(next http.Handler, jwtKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const logErrPrefix = "middleware.AuthorizationRequired():"

		authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if authToken == "" {
			cookie, err := r.Cookie("authToken")
			if err != nil {
				httperrorwriter.WriteError(w, appErrors.ErrNoTokenProvided, http.StatusUnauthorized, logErrPrefix)
				return
			}

			authToken = cookie.Value
		}

		_, err := jwt.CheckIsAdminInJWT(authToken, jwtKey)
		if err != nil {
			if errors.Is(err, appErrors.ErrTokenIsInvalid) {
				httperrorwriter.WriteError(w, appErrors.ErrTokenIsInvalid, http.StatusUnauthorized, logErrPrefix)
				return
			}

			logger.Logger().Errorln(logErrPrefix, zap.Error(err))
			httperrorwriter.WriteError(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError, logErrPrefix)
			return
		}

		next.ServeHTTP(w, r)
	})
}
