package main

import (
	"fmt"
	"github.com/bvankampen/kubeconfig-selector/pkg/app"
	"os"
)

var (
	Version  = "0"
	CommitId = "0"
)

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			{
				fmt.Printf("Kubernetes Kubeconfig Cluster Selector\n"+
					"(C) 2023 Bas van Kampen <bas@ping6.nl>\n"+
					"Version %s-%s\n", Version, CommitId)
			}
		}
	} else {
		application := app.Application{}
		application.Run()
	}

}
