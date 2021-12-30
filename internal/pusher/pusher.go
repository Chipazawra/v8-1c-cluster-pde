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

func WithConfig(config PusherConfig) PusherOption {
	return func(p *Pusher) {
		p.url = fmt.Sprintf("%s:%s", config.PUSH_HOST, config.PUSH_PORT)
		p.intervalMillis = config.PUSH_INTERVAL
	}
}

func New(collector prometheus.Collector, opts ...PusherOption) *Pusher {

	p := &Pusher{
		collector: collector,
		jobName:   defaultJobname,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.pusher = push.New(p.url, p.jobName).Collector(collector)

	return p
}

func (p *Pusher) Run(ctx context.Context, errchan chan<- error) {
	log.Printf("v8-1c-cluster-pde: pusher %v", p.url)
	ticker := time.NewTicker(time.Duration(p.intervalMillis * int(time.Microsecond)))
Loop:
	for {
		select {
		case <-ticker.C:
			err := p.pusher.Push()
			if err != nil {
				errchan <- fmt.Errorf("puser: %v", err)
				break Loop
			}
		case <-ctx.Done():
			log.Println("INFO: pusher context done")
			break Loop
		}
	}
}
