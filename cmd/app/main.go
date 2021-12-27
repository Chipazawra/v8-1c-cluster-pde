package main

import (
	"context"
	"log"

	rascli "github.com/khorevaa/ras-client"
	"github.com/khorevaa/ras-client/serialize"
)

func main() {

	ctx := context.Background()

	client := rascli.NewClient("192.168.10.233:1545")
	clusters, err := client.GetClusters(ctx)

	if err != nil {
		log.Printf("rac-1c-impl: %s", err)
	}

	for _, cls := range clusters {

		log.Printf("%v", cls)

		workingProcesses, err := client.GetWorkingProcesses(ctx, cls.UUID)

		if err != nil {
			log.Printf("rac-1c-impl: %s", err)
		}

		workingProcesses.Each(func(info *serialize.ProcessInfo) {
			log.Printf("ProcessInfo: %#v", info)
		})

	}

}
