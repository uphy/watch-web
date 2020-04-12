package main

import (
	"os"

	"github.com/uphy/watch-web/cli"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	app := cli.New()
	return app.Run(os.Args)
}
