package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pusher "github.com/Chipazawra/v8-1c-cluster-pde/internal/pusher"
	"github.com/Chipazawra/v8-1c-cluster-pde/internal/rpHostsCollector"
	"github.com/caarlos0/env"
	rascli "github.com/khorevaa/ras-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	push string = "push"
	pull string = "pull"
)

var (
	conf          Config
	RAS_HOST      string
	RAS_PORT      string
	PULL_EXPOSE   string
	MODE          string
	PUSH_INTERVAL int
	PUSH_HOST     string
	PUSH_PORT     string
)

func init() {

	if err := env.Parse(&conf); err != nil {
		log.Fatalf("app: config...")
	}

	log.Printf("v8-1c-cluster-pde: config from env:\n %#v", conf)

	flag.StringVar(&RAS_HOST, "ras-host", "", "cluster host.")
	flag.StringVar(&RAS_PORT, "ras-port", "", "cluster port.")
	flag.StringVar(&PULL_EXPOSE, "pull-expose", "", "metrics port.")
	flag.StringVar(&MODE, "mode", "", "mode push or pull")
	flag.IntVar(&PUSH_INTERVAL, "push-interval", 0, "mode push or pull")
	flag.StringVar(&PUSH_HOST, "push-host", "", "pushgateway host")
	flag.StringVar(&PUSH_PORT, "push-port", "", "pushgateway port")
	flag.Parse()

	if RAS_HOST != "" {
		conf.RAS_HOST = RAS_HOST
	}

	if RAS_PORT != "" {
		conf.RAS_PORT = RAS_PORT
	}

	if PULL_EXPOSE != "" {
		conf.PULL_EXPOSE = PULL_EXPOSE
	}

	if MODE != "" {
		conf.MODE = MODE
	}

	if PUSH_INTERVAL != 0 {
		conf.PUSH_INTERVAL = PUSH_INTERVAL
	}

	if PUSH_HOST != "" {
		conf.PUSH_HOST = PUSH_HOST
	}

	if PUSH_PORT != "" {
		conf.PUSH_PORT = PUSH_PORT
	}

	log.Printf("v8-1c-cluster-pde: overrided config from stdin:\n%#v", conf)
}

func Run() error {

	rcli := rascli.NewClient(fmt.Sprintf("%s:%s", conf.RAS_HOST, conf.RAS_PORT))
	rcli.AuthenticateAgent(conf.CLS_USER, conf.CLS_PASS)
	defer rcli.Close()

	log.Printf("v8-1c-cluster-pde: connected to RAS %v",
		fmt.Sprintf("%s:%s", conf.RAS_HOST, conf.RAS_PORT),
	)

	ctx, cancel := context.WithCancel(context.Background())

	sigchan := make(chan os.Signal, 1)
	defer close(sigchan)
	errchan := make(chan error, 1)
	defer close(errchan)

	signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	switch conf.MODE {
	case push:
		go RunPusher_(ctx, errchan, rcli)
	case pull:
		go RunPuller_(ctx, errchan, rcli)
	default:
		{
			cancel()
			return fmt.Errorf("v8-1c-cluster-pde: %v", "undefined mode")
		}
	}

	log.Printf("v8-1c-cluster-pde: received signal %v",
		<-sigchan,
	)
	cancel()

	return <-errchan
}

func RunPuller(rasapi rascli.Api) error {

	log.Printf("v8-1c-cluster-pde: runing in %v mode", conf.MODE)
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(rpHostsCollector.New(rasapi))

	http.Handle("/metrics",
		promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	)

	log.Printf("v8-1c-cluster-pde: listen %v", fmt.Sprintf("%s:%s", "", conf.PULL_EXPOSE))

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", "", conf.PULL_EXPOSE), nil)
	if err != nil {
		return fmt.Errorf("app: %v", err)
	}

	return nil
}

func RunPusher(rasapi rascli.Api) error {
	log.Printf("v8-1c-cluster-pde: runing in %v mode pushgateway %v\n",
		conf.MODE, fmt.Sprintf("%s:%s", conf.PUSH_HOST, conf.PUSH_PORT))
	return pusher.New(
		rpHostsCollector.New(rasapi),
		fmt.Sprintf("%s:%s", conf.PUSH_HOST, conf.PUSH_PORT),
		pusher.WithInterval(500),
	).Run(context.Background())
}

func RunPuller_(ctx context.Context, errchan chan<- error, rasapi rascli.Api) {

	log.Printf("v8-1c-cluster-pde: runing in %v mode", conf.MODE)
	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(rpHostsCollector.New(rasapi))

	mux := http.NewServeMux()
	mux.Handle("/metrics",
		promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	)
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", "", conf.PULL_EXPOSE),
		Handler: mux,
	}

	go func() {
		errchan <- srv.ListenAndServe()
	}()
	log.Printf("v8-1c-cluster-pde: listen %v", fmt.Sprintf("%s:%s", "", conf.PULL_EXPOSE))

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		errchan <- fmt.Errorf("v8-1c-cluster-pde: server shutdown with err: %v", err)
	}
	log.Printf("v8-1c-cluster-pde: server shutdown")
}

func RunPusher_(ctx context.Context, errchan chan<- error, rasapi rascli.Api) {
	log.Printf("v8-1c-cluster-pde: runing in %v mode pushgateway %v\n",
		conf.MODE, fmt.Sprintf("%s:%s", conf.PUSH_HOST, conf.PUSH_PORT))

	errchan <- pusher.New(
		rpHostsCollector.New(rasapi),
		fmt.Sprintf("%s:%s", conf.PUSH_HOST, conf.PUSH_PORT),
		pusher.WithInterval(500),
	).Run(ctx)
}
