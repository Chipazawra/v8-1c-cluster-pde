package pusher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	defaultIntervalMillis int    = 500
	defaultJobname        string = "v8-1C-cluster-pde"
)

type Pusher struct {
	intervalMillis int
	collector      prometheus.Collector
	url            string
	jobName        string
	pusher         *push.Pusher
}

type PusherOption func(*Pusher)

func WithInterval(millis int) PusherOption {
	return func(p *Pusher) {
		p.intervalMillis = millis
	}
}

func WithJobName(Name string) PusherOption {
	return func(p *Pusher) {
		p.jobName = Name
	}
}

func New(collector prometheus.Collector, url string, opts ...PusherOption) *Pusher {

	p := &Pusher{
		collector:      collector,
		intervalMillis: defaultIntervalMillis,
		url:            url,
		jobName:        defaultJobname,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.pusher = push.New(url, p.jobName).Collector(collector)

	return p
}

func (p *Pusher) Run(ctx context.Context) error {

	ticker := time.NewTicker(time.Duration(p.intervalMillis * int(time.Microsecond)))
	done := make(chan error)
	go func(done chan error) {
	Loop:
		for {
			select {
			case <-ticker.C:
				err := p.pusher.Push()
				if err != nil {
					done <- fmt.Errorf("puser: %v", err)
					break Loop
				}
			case <-ctx.Done():
				log.Println("INFO: pusher context complete")
				done <- nil
				break Loop
			}
		}
		close(done)
	}(done)

	return <-done
}
