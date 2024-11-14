package main

import (
	"fmt"
	"log"

	"github.com/Chipazawra/v8-1c-cluster-pde/internal/app"
	"github.com/kardianos/service"
)

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	serviceConfig := &service.Config{
		Name:        "V8_RAC_CLUSTER_PDE",
		DisplayName: "V8 1C CLUSTER PDE",
		Description: "Prometheus 1C cluster exporter",
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}

}
