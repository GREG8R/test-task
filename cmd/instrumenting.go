package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"

	"btcn_srv/pkg/services/bitcoin_service"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           bitcoin_service.BitcoinService
}

func (mw instrumentingMiddleware) SendMoney(ctx context.Context, request interface{}) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SendMoney", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.SendMoney(ctx, request)
	return
}

func (mw instrumentingMiddleware) GetHistory(ctx context.Context, request interface{}) (result []bitcoin_service.GetHistoryResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetHistory", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result, err = mw.next.GetHistory(ctx, request)
	return
}
