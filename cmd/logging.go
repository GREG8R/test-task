package main

import (
	"context"
	"fmt"
	"time"

	"btcn_srv/pkg/services/bitcoin_service"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   bitcoin_service.BitcoinService
}

func (mw loggingMiddleware) SendMoney(ctx context.Context, request interface{}) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SendMoney",
			"input", fmt.Sprint(request),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.SendMoney(ctx, request)
	return
}

func (mw loggingMiddleware) GetHistory(ctx context.Context, request interface{}) (resp []bitcoin_service.GetHistoryResponse, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetHistory",
			"input", fmt.Sprint(request),
			"response", fmt.Sprint(resp),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	resp, err = mw.next.GetHistory(ctx, request)
	return
}
