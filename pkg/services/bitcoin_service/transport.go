package bitcoin_service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func MakeHTTPHandler(svc BitcoinService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /send
	// GET     /history
	r.Methods("POST").Path("/send").Handler(httptransport.NewServer(
		MakeSendMoneyEndpoint(svc),
		DecodeSendMoneyRequest,
		EncodeResponse,
		options...,
	))
	r.Methods("GET").Path("/history").Handler(httptransport.NewServer(
		MakeGetHistoryEndpoint(svc),
		DecodeGetHistoryRequest,
		EncodeResponse,
		options...,
	))

	return r
}

func MakeSendMoneyEndpoint(svc BitcoinService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(SendMoneyRequest)

		err := sendMoneyValidation(req)
		if err != nil {
			return nil, err
		}

		v, err := svc.SendMoney(req)
		if err != nil {
			return nil, err
		}

		return SendMoneyResponse{v, ""}, nil
	}
}

func MakeGetHistoryEndpoint(svc BitcoinService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(GetHistoryRequest)

		err := getHistoryValidation(req)
		if err != nil {
			return nil, err
		}

		response := svc.GetHistory(req)
		return response, nil
	}
}

func DecodeSendMoneyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SendMoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrErrorParseRequest
	}
	return request, nil
}

func DecodeGetHistoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GetHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrErrorParseRequest
	}
	return request, nil
}

type errorer interface {
	error() error
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

type SendMoneyRequest struct {
	Amount decimal.Decimal `json:"amount"`
	Date   time.Time       `json:"datetime"`
}

type SendMoneyResponse struct {
	Body string `json:"body,omitempty"`
	Err  string `json:"err,omitempty"`
}

type GetHistoryRequest struct {
	StartDate time.Time `json:"startDatetime"`
	EndDate   time.Time `json:"endDatetime"`
}

type GetHistoryResponse struct {
	Amount decimal.Decimal `json:"amount"`
	Date   time.Time       `json:"datetime"`
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var (
	ErrErrorParseRequest = errors.New("can't parse request")
	ErrNegativeAmount    = errors.New("can't get negative amount")
	ErrStartDateLater    = errors.New("start date should be earlier than end date")
)

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

func sendMoneyValidation(req SendMoneyRequest) error {
	if req.Amount.IsNegative() {
		return ErrNegativeAmount
	}
	return nil
}

func getHistoryValidation(req GetHistoryRequest) error {
	if req.EndDate.UTC().Nanosecond() < req.StartDate.UTC().Nanosecond() {
		return ErrStartDateLater
	}
	return nil
}
