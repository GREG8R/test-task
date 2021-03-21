package bitcoin_service

func sendMoneyValidation(req SendMoneyRequest) error {
	if req.Amount.IsNegative() {
		return ErrNegativeAmount
	}
	return nil
}

func getHistoryValidation(req GetHistoryRequest) error {
	if req.EndDate.Before(req.StartDate) {
		return ErrStartDateLater
	}
	return nil
}
