package main

import "github.com/urfave/cli/v2"

type trafficType int

const (
	trafficTypeHttp trafficType = iota
)

func trafficTypeFromString(s string) trafficType {
	switch s {
	case "http":
		return trafficTypeHttp
	}
	return -1
}

type req struct {
	ctx *cli.Context
}

func newReq(ctx *cli.Context) *req {
	return &req{ctx: ctx}
}
