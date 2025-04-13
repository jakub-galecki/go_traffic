package main

import "fmt"

type httpHandler struct{}

func (h *httpHandler) do(r *req) error {
	fmt.Println("running http traffic generator")

	return nil
}
