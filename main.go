package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

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

	r := newReq(c)
	return g.handlers[t].do(r)
}

func main() {
	gen := &trafficGenerator{
		handlers: map[trafficType]handler{
			trafficTypeHttp: &httpHandler{},
		},
	}

	app := &cli.App{
		Name:  "Go Traffic Generator",
		Usage: "Generate http traffic for stress testing applications",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Value:   "http",
				Usage:   "type of traffic that should be generated: http",
				Aliases: []string{"t"},
			},
		},
		Action: gen.run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
