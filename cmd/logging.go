package main

import (
	"time"

	"btcn_srv/pkg/services/bitcoin_service"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   bitcoin_service.BitcoinService
}

func (mw loggingMiddleware) SendMoney(s interface{}) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SendMoney",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.SendMoney(s)
	return
}

func (mw loggingMiddleware) GetHistory(s interface{}) (n []bitcoin_service.GetHistoryResponse) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetHistory",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.next.GetHistory(s)
	return
}
