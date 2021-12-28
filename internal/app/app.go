package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Chipazawra/v8-1c-cluster-pde/internal/rpHostsCollector"
	"github.com/caarlos0/env"
	rascli "github.com/khorevaa/ras-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	conf       Config
	hostFlag   string
	portFlag   string
	exposeFlag string
	ctx        context.Context
)

func init() {
	if err := env.Parse(&conf); err != nil {
		log.Fatalf("app: config...")
	}

	flag.StringVar(&hostFlag, "host", "", "cluster host.")
	flag.StringVar(&portFlag, "port", "", "cluster port.")
	flag.StringVar(&exposeFlag, "expose", "", "metrics port.")
	flag.Parse()

	if hostFlag != "" {
		conf.Host = hostFlag
	}

	if portFlag != "" {
		conf.Port = portFlag
	}

	if exposeFlag != "" {
		conf.Expose = exposeFlag
	}

	ctx = context.Background()

}

func Run() error {

	rcli := rascli.NewClient(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	rcli.AuthenticateAgent(conf.User, conf.Pass)
	defer rcli.Close()

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(rpHostsCollector.New(rcli))

	http.Handle("/metrics",
		promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", "", conf.Expose), nil)
	if err != nil {
		return Errorf(err)
	}

	return nil
}

func Errorf(err error) error {
	return fmt.Errorf("app: %v", err)
}
