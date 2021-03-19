package bitcoin_service

import (
	"fmt"
)

// StringService provides operations on strings.
type BitcoinService interface {
	SendMoney(interface{}) (string, error)
	GetHistory(interface{}) []GetHistoryResponse
}

type BtcnService struct{}

func (BtcnService) SendMoney(s interface{}) (string, error) {
	req := s.(SendMoneyRequest)

	fmt.Print(req)

	return "", nil
}

func (BtcnService) GetHistory(s interface{}) []GetHistoryResponse {

	return nil
}
