package bitcoin_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ErrErrorParseRequest = errors.New("can't parse request")
	ErrNegativeAmount    = errors.New("can't get negative amount")
	ErrStartDateLater    = errors.New("start date should be earlier than end date")
)

type errorer interface {
	error() error
}

func MakeHTTPHandler(svc BitcoinService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST /send
	// example:
	//curl --location --request POST 'localhost:8080/send' \
	//--header 'Content-Type: application/json' \
	//--data-raw '{
	//    "datetime": "2019-10-05T14:48:01+01:00",
	//    "amount": 1.2
	//}'
	r.Methods("POST").Path("/send").Handler(httptransport.NewServer(
		MakeSendMoneyEndpoint(svc),
		DecodeSendMoneyRequest,
		EncodeResponse,
		options...,
	))

	// GET /history
	// example:
	//curl --location --request GET 'localhost:8080/history' \
	//--header 'Content-Type: application/json' \
	//--data-raw '{
	//    "startDatetime": "2019-10-05T13:48:01+01:00",
	//    "endDatetime": "2019-10-05T15:48:01+01:00"
	//}
	//'
	r.Methods("GET").Path("/history").Handler(httptransport.NewServer(
		MakeGetHistoryEndpoint(svc),
		DecodeGetHistoryRequest,
		EncodeResponse,
		options...,
	))

	// GET /metrict
	// example:
	// curl --location --request GET 'localhost:8080/metrics'
	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	return r
}

func MakeSendMoneyEndpoint(svc BitcoinService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendMoneyRequest)
		err := svc.SendMoney(ctx, req)
		return nil, err
	}
}

func MakeGetHistoryEndpoint(svc BitcoinService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetHistoryRequest)
		response, err := svc.GetHistory(ctx, req)
		return response, err
	}
}

func DecodeSendMoneyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SendMoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrErrorParseRequest
	}

	err := sendMoneyValidation(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeGetHistoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GetHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrErrorParseRequest
	}

	err := getHistoryValidation(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrErrorParseRequest:
		return http.StatusBadRequest
	case ErrNegativeAmount:
		return http.StatusBadRequest
	case ErrStartDateLater:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
