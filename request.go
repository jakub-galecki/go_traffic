package main

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
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
		method  string
		headers http.Header
	}
}

func newReq(ctx *cli.Context, t trafficType) *req {
	parseHeader := func(raw string) (string, string, bool) {
		splitted := strings.Split(raw, ":")
		if len(splitted) < 2 {
			return "", "", false
		}
		return splitted[0], splitted[1], true
	}

	r := &req{ctx: ctx}
	r.parsed.target = ctx.String("url")
	r.parsed.timeout = ctx.Duration("timeout")
	r.parsed.count = ctx.Int("size")
	r.parsed.workers = ctx.Int("workers")
	r.parsed.tt = t
	r.parsed.headers = make(http.Header)
	r.parsed.method = ctx.String("method")

	for _, h := range ctx.StringSlice("headers") {
		key, value, ok := parseHeader(h)
		if !ok {
			continue
		}
		r.parsed.headers.Set(key, value)
	}

	slog.Debug("[newReq] created request",
		slog.String("url", r.parsed.target),
		slog.Duration("timeout", r.parsed.timeout),
		slog.Int("count", r.parsed.count),
		slog.Int("workers", r.parsed.workers))

	return r
}
