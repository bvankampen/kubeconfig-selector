package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bvankampen/kubeconfig-selector/internal/selector"
	"github.com/bvankampen/kubeconfig-selector/internal/ui"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var (
	Version  = "0.0"
	CommitId = "dev"
)

func main() {
	cmd := &cli.Command{
		Name:    "ks",
		Version: fmt.Sprintf("%s (%s)", Version, CommitId),
		Usage:   "Select kubeconfig",
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return ctx, nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug",
			},
		},
		Action: run,
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	s, err := selector.New(cmd)
	if err != nil {
		return err
	}

	var ui ui.UI
	if err := ui.Init(s.Cmd(), s.AppConfig(), s.KubeConfigs(), s.ActiveConfig(), s.Debug()); err != nil {
		return err
	}
	return ui.Run()
}
