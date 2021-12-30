package puller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Puller struct {
	collector prometheus.Collector
	expose    string
}

type PullerOption func(*Puller)

func WithConfig(config PullerConfig) PullerOption {
	return func(p *Puller) {
		p.expose = config.PULL_EXPOSE

	}
}

func New(collector prometheus.Collector, opts ...PullerOption) *Puller {
	p := &Puller{
		collector: collector,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Puller) Run(ctx context.Context, errchan chan<- error) {

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(p.collector)

	mux := http.NewServeMux()
	mux.Handle("/metrics",
		promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	)
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", "", p.expose),
		Handler: mux,
	}

	go func() {
		errchan <- srv.ListenAndServe()
	}()
	log.Printf("v8-1c-cluster-pde: puller listen %v", fmt.Sprintf("%s:%s", "", p.expose))

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		errchan <- fmt.Errorf("v8-1c-cluster-pde: puller server shutdown with err: %v", err)
	}
	log.Printf("v8-1c-cluster-pde: puller server shutdown")
}
