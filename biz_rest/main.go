package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	svc := NewBizService()

	endpoint := MakeBizEndpoint(svc)

	r := MakeHttpHandler(ctx, endpoint, logger)

	go func() {
		fmt.Println("http server start at port:8000")
		handler := r
		errChan <- http.ListenAndServe(":8000", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("http.ListenAndServe error:%v\n", <-errChan)
}
