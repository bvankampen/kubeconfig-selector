package main

import (
	"fmt"
	"github.com/bvankampen/kubeconfig-selector/internal/selector"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	Version  = "0.0"
	CommitId = "dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "cluster"
	app.Version = fmt.Sprintf("%s (%s)", Version, CommitId)
	app.Usage = "Select kubeconfig"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debugf("Loglevel set to [%v]", logrus.DebugLevel)
		}
		return nil
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug",
		}}
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx *cli.Context) {
	s, err := selector.New(*ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}
