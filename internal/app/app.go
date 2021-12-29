package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
)

func init() {
	if err := env.Parse(&conf); err != nil {
		log.Fatalf("app: config...")
	}

	log.SetOutput(os.Stdout)

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
}

func Run() error {

	rcli := rascli.NewClient(fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	rcli.AuthenticateAgent(conf.User, conf.Pass)
	log.Printf("cluster-pde connected to RAS: %v", fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	defer rcli.Close()

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(rpHostsCollector.New(rcli))

	http.Handle("/metrics",
		promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	)

	log.Printf("cluster-pde is running on: %v", fmt.Sprintf("%s:%s", "", conf.Expose))

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", "", conf.Expose), nil)
	if err != nil {
		return fmt.Errorf("app: %v", err)
	}

	return nil
}
