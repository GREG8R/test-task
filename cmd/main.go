package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	//stdprometheus "github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"

	"btcn_srv/pkg/services/bitcoin_service"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	//fieldKeys := []string{"method", "error"}
	//requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	//	Namespace: "my_group",
	//	Subsystem: "string_service",
	//	Name:      "request_count",
	//	Help:      "Number of requests received.",
	//}, fieldKeys)
	//requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	//	Namespace: "my_group",
	//	Subsystem: "string_service",
	//	Name:      "request_latency_microseconds",
	//	Help:      "Total duration of requests in microseconds.",
	//}, fieldKeys)
	//countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	//	Namespace: "my_group",
	//	Subsystem: "string_service",
	//	Name:      "count_result",
	//	Help:      "The result of each count method.",
	//}, []string{}) // no fields here

	var svc bitcoin_service.BitcoinService
	svc = bitcoin_service.BtcnService{}
	svc = loggingMiddleware{logger, svc}
	//svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

	h := bitcoin_service.MakeHTTPHandler(svc, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)

	//http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
