package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/cli"
)

func main() {
	log := logrus.New()
	log.SetFormatter(new(logrus.TextFormatter))

	if err := run(log); err != nil {
		panic(err)
	}
}

func run(log *logrus.Logger) error {
	app := cli.New(log)
	return app.Run(os.Args)
}
