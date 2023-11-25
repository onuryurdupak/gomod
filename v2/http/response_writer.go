package gmhttp

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
}

// NewResponseWriter creates an utility instance for handling of http responses.
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

// WriteCustomJsonResponse serializes input res and creates response payload from it.
func (hu *ResponseWriter) WriteCustomJsonResponse(w http.ResponseWriter, statusCode int, res interface{}) (writtenRes []byte, err error) {
	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(resJson)
	if err != nil {
		return nil, err
	}

	return resJson, nil
}
