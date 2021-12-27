package v8clustercollector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type clusterCollector struct {
	rpHosts *prometheus.Desc
}

func NewClusterCollector() prometheus.Collector {
	return &clusterCollector{
		rpHosts: prometheus.NewDesc(
			"rp_hosts",
			"count of active rp hosts on cluster",
			nil, nil),
	}
}

func (c *clusterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.rpHosts
}

func (c *clusterCollector) Collect(ch chan<- prometheus.Metric) {

}
