package jsonhttpvalidator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	httperrorwriter "github.com/PoorMercymain/filmoteka/pkg/http-error-writer"
	jsonduplicatechecker "github.com/PoorMercymain/filmoteka/pkg/json-duplicate-checker"
	jsonmimechecker "github.com/PoorMercymain/filmoteka/pkg/json-mime-checker"
)

func ValidateJSONRequest(w http.ResponseWriter, r *http.Request, logErrPrefix string) error {
	if !jsonmimechecker.IsJSONContentTypeCorrect(r) {
		httperrorwriter.WriteError(w, appErrors.ErrWrongMIME, http.StatusBadRequest, logErrPrefix)
		return appErrors.ErrWrongRequestWithJSON
	}

	bytesToCheck, err := io.ReadAll(r.Body)
	if err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return appErrors.ErrWrongRequestWithJSON
	}

	reader := bytes.NewReader(bytes.Clone(bytesToCheck))

	err = jsonduplicatechecker.CheckDuplicatesInJSON(json.NewDecoder(reader), nil)
	if err != nil {
		httperrorwriter.WriteError(w, err, http.StatusBadRequest, logErrPrefix)
		return appErrors.ErrWrongRequestWithJSON
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bytesToCheck))

	return nil
}
