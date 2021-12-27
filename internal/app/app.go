package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	rclient "github.com/khorevaa/ras-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	conf     Config
	hostFlag string
	portFlag string
)

func init() {
	if err := env.Parse(&conf); err != nil {
		log.Fatalf("app: config...")
	}

	flag.StringVar(&hostFlag, "host", "", "cluster host.")
	flag.StringVar(&portFlag, "port", "", "cluster port.")
	flag.Parse()

	if hostFlag != "" {
		conf.Host = hostFlag
	}

	if portFlag != "" {
		conf.Port = portFlag
	}
}

func Run() error {

	_ = rclient.NewClient(fmt.Sprintf("%s:%s", conf.Host, conf.Port))

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister()
	_ = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "",
			Help:        "",
			ConstLabels: map[string]string{},
		})

	_ = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "",
			Help:        "",
			ConstLabels: map[string]string{},
		},
		nil,
	)

	http.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))

	return fmt.Errorf("app: not implemented")

}
