package bitcoin_service

import (
	"context"

	"btcn_srv/pkg/pg_storage"
)

type BitcoinService interface {
	SendMoney(ctx context.Context, request interface{}) error
	GetHistory(ctx context.Context, request interface{}) ([]GetHistoryResponse, error)
}

type BtcnService struct {
	Storage pg_storage.PgStorage
}

func (s BtcnService) SendMoney(ctx context.Context, request interface{}) error {
	req := request.(SendMoneyRequest)

	// save all dates like UTC
	err := s.Storage.SaveMoney(ctx, req.Amount, req.Date.UTC())

	return err
}

func (s BtcnService) GetHistory(ctx context.Context, request interface{}) ([]GetHistoryResponse, error) {
	var response []GetHistoryResponse
	req := request.(GetHistoryRequest)

	// get all dates by UTC
	result, err := s.Storage.GetHistory(ctx, req.StartDate.UTC(), req.EndDate.UTC())
	if err != nil {
		return nil, err
	}

	// get location by date from request
	loc := req.StartDate.Location()
	for _, r := range result {
		response = append(response, GetHistoryResponse{
			Amount: r.Amount,
			Date:   r.Hour.In(loc),
		})
	}

	return response, nil
}
