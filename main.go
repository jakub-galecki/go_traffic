package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	_ = slog.SetLogLoggerLevel(slog.LevelDebug)
}

type handler interface {
	do(*req) error
}

type trafficGenerator struct {
	handlers map[trafficType]handler
}

func (g *trafficGenerator) run(c *cli.Context) error {
	t := trafficTypeFromString(c.String("type"))
	if t < 0 {
		return fmt.Errorf("invalid traffic type: %s", c.String("type"))
	}

	if _, ok := g.handlers[t]; !ok {
		return fmt.Errorf("handler for traffic type not found: %s", c.String("type"))
	}

	r := newReq(c, t)
	return g.handlers[t].do(r)
}

func main() {
	gen := &trafficGenerator{
		handlers: map[trafficType]handler{
			trafficTypeHttp: newHttpHandler(),
		},
	}

	app := &cli.App{
		Name:  "Go Traffic Generator",
		Usage: "Generate http traffic for stress testing applications",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Required: true,
				Usage:    "url of the target to be tested",
				Aliases:  []string{"u"},
			},
			&cli.StringFlag{
				Name:    "type",
				Value:   "http",
				Usage:   "type of traffic that should be generated: http",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:  "timeout",
				Value: "500ms",
				Usage: "for how long traffic should be generated",
			},
			&cli.StringFlag{
				Name:  "method",
				Value: "GET",
				Usage: "what type of the http should be sent: GET,POST",
			},
			&cli.Int64Flag{
				Name:    "size",
				Value:   10,
				Usage:   "how many requests should be generated",
				Aliases: []string{"s"},
			},
			&cli.Int64Flag{
				Name:    "workers",
				Value:   10,
				Usage:   "number of concurrent workers that will generate request traffic",
				Aliases: []string{"w"},
			},
			&cli.StringSliceFlag{
				Name:    "headers",
				Usage:   "headers that will be passed on the request, with format key:value",
				Aliases: []string{"h"},
			},
		},
		Action: gen.run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
