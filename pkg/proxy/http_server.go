package proxy

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoint Endpoints) http.Handler {

	m := mux.NewRouter()

	m.Methods("POST").Path("/multiply/").Handler(httptransport.NewServer(
		endpoint.MultiplyEndpoint,
		decodeHTTPMultiplyRequest,
		encodeHTTPGenericResponse,
	))

	return m
}

func decodeHTTPMultiplyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req = MultiplyRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	// switch err {
	// case addservice.ErrTwoZeroes, addservice.ErrMaxSizeExceeded, addservice.ErrIntOverflow:
	// 	return http.StatusBadRequest
	// }
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
