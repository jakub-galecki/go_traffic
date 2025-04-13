package main

import (
	"github.com/urfave/cli/v2"
	"log/slog"
	"time"
)

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

	parsed struct {
		target  string
		tt      trafficType
		timeout time.Duration
		count   int
		workers int
	}
}

func newReq(ctx *cli.Context, t trafficType) *req {
	r := &req{ctx: ctx}
	r.parsed.target = ctx.String("url")
	r.parsed.timeout = ctx.Duration("timeout")
	r.parsed.count = ctx.Int("size")
	r.parsed.workers = ctx.Int("workers")
	r.parsed.tt = t

	slog.Debug("[newReq] created request",
		slog.String("url", r.parsed.target),
		slog.Duration("timeout", r.parsed.timeout),
		slog.Int("count", r.parsed.count),
		slog.Int("workers", r.parsed.workers))

	return r
}
