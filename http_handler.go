package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	_maxWorkers = 10
)

type httpHandler struct {
	client *http.Client
}

func newHttpHandler() *httpHandler {
	return &httpHandler{
		client: http.DefaultClient,
	}
}

func requestGenerator(r *req) func() (*http.Request, error) {
	return func() (*http.Request, error) {
		rr, err := http.NewRequest(r.parsed.method, r.parsed.target, nil)
		if err != nil {
			return nil, err
		}
		return rr, nil
	}
}

func (h *httpHandler) doInternal(r *http.Request) error {
	start := time.Now()
	res, err := h.client.Do(r)
	_ = time.Since(start) // elapsed
	if err != nil {
		return err
	}
	// parse res
	_ = res
	return nil
}

func (h *httpHandler) do(r *req) error {
	fmt.Println("running http traffic generator")

	var (
		total = atomic.Int64{}
		gen   = requestGenerator(r)
	)

	total.Store(int64(r.parsed.count))

	for total.Load() > 0 {
		rr, err := gen()
		if err != nil {
			return err
		}
		errg := &errgroup.Group{}
		errg.SetLimit(_maxWorkers)
		for i := 0; i < r.parsed.workers; i++ {
			errg.Go(func() error {
				return h.doInternal(rr)
			})
		}
	}
	return nil
}
