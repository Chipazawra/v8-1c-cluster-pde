package rpHostsCollector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	rascli "github.com/khorevaa/ras-client"
	"github.com/khorevaa/ras-client/serialize"
	"github.com/prometheus/client_golang/prometheus"
)

type rpHostsCollector struct {
	ctx                 context.Context
	wg                  sync.WaitGroup
	rasapi              rascli.Api
	rpHosts             *prometheus.Desc
	Duration            *prometheus.Desc
	MemorySize          *prometheus.Desc
	Connections         *prometheus.Desc
	AvgThreads          *prometheus.Desc
	AvailablePerfomance *prometheus.Desc
	Capacity            *prometheus.Desc
	MemoryExcessTime    *prometheus.Desc
	SelectionSize       *prometheus.Desc
	AvgBackCallTime     *prometheus.Desc
	AvgCallTime         *prometheus.Desc
	AvgDbCallTime       *prometheus.Desc
	AvgLockCallTime     *prometheus.Desc
	AvgServerCallTime   *prometheus.Desc
	Running             *prometheus.Desc
	Enable              *prometheus.Desc
}

func New(rasapi rascli.Api) prometheus.Collector {

	proccesLabels := []string{"cluster", "pid", "host", "port", "startedAt"}

	return &rpHostsCollector{
		ctx:    context.Background(),
		rasapi: rasapi,
		rpHosts: prometheus.NewDesc(
			"rp_hosts_active",
			"count of active rp hosts on cluster",
			[]string{"cluster"}, nil),
		MemorySize: prometheus.NewDesc(
			"rp_hosts_memory",
			"count of active rp hosts on cluster",
			proccesLabels, nil),
		Duration: prometheus.NewDesc(
			"rp_hosts_scrape_duration",
			"the time in milliseconds it took to collect the metrics",
			nil, nil),
		Connections: prometheus.NewDesc(
			"rp_hosts_connections",
			"number of connections to host",
			proccesLabels, nil),
		AvgThreads: prometheus.NewDesc(
			"rp_hosts_avg_threads",
			"average number of client threads",
			proccesLabels, nil),
		AvailablePerfomance: prometheus.NewDesc(
			"rp_hosts_available_perfomance",
			"available host performance",
			proccesLabels, nil),
		Capacity: prometheus.NewDesc(
			"rp_hosts_capacity",
			"host capacity",
			proccesLabels, nil),
		MemoryExcessTime: prometheus.NewDesc(
			"rp_hosts_memory_excess_time",
			"host memory excess time",
			proccesLabels, nil),
		SelectionSize: prometheus.NewDesc(
			"rp_hosts_selection_size",
			"host selection size",
			proccesLabels, nil),
		AvgBackCallTime: prometheus.NewDesc(
			"rp_hosts_avg_back_call_time",
			"host avg back call time",
			proccesLabels, nil),
		AvgCallTime: prometheus.NewDesc(
			"rp_hosts_avg_call_time",
			"host avg call time",
			proccesLabels, nil),
		AvgDbCallTime: prometheus.NewDesc(
			"rp_hosts_avg_db_call_time",
			"host avg db call time",
			proccesLabels, nil),
		AvgLockCallTime: prometheus.NewDesc(
			"rp_hosts_avg_lock_call_time",
			"host avg lock call time",
			proccesLabels, nil),
		AvgServerCallTime: prometheus.NewDesc(
			"rp_hosts_avg_server_call_time",
			"host avg server call time",
			proccesLabels, nil),
		Enable: prometheus.NewDesc(
			"rp_hosts_enable",
			"host enable",
			proccesLabels, nil),
		Running: prometheus.NewDesc(
			"rp_hosts_running",
			"host enable",
			proccesLabels, nil),
	}
}

func (c *rpHostsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.rpHosts
}

func (c *rpHostsCollector) Collect(ch chan<- prometheus.Metric) {

	start := time.Now()

	сlusters, err := c.rasapi.GetClusters(c.ctx)
	if err != nil {
		log.Printf("rpHostsCollector: %v", err)
	}

	for _, сluster := range сlusters {
		c.wg.Add(1)
		go c.funInCollect(ch, *сluster)
	}
	c.wg.Wait()

	ch <- prometheus.MustNewConstMetric(
		c.Duration,
		prometheus.GaugeValue,
		float64(time.Since(start).Milliseconds()))
}

func (c *rpHostsCollector) funInCollect(ch chan<- prometheus.Metric, clusterInfo serialize.ClusterInfo) {

	var (
		rpHostsCount int
	)

	workingProcesses, err := c.rasapi.GetWorkingProcesses(c.ctx, clusterInfo.UUID)
	if err != nil {
		log.Printf("rpHostsCollector: %v", err)
	}

	workingProcesses.Each(func(proccesInfo *serialize.ProcessInfo) {

		var (
			proccesLabelsVal []string = []string{
				clusterInfo.Name,
				proccesInfo.Pid,
				fmt.Sprint(proccesInfo.Host),
				fmt.Sprint(proccesInfo.Port),
				// lag MSK+3
				proccesInfo.StartedAt.In(time.UTC).Format("2006-01-02 15:04:05"),
			}
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemorySize,
			prometheus.GaugeValue,
			float64(proccesInfo.MemorySize),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Connections,
			prometheus.GaugeValue,
			float64(proccesInfo.Connections),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvgThreads,
			prometheus.GaugeValue,
			float64(proccesInfo.AvgThreads),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvailablePerfomance,
			prometheus.GaugeValue,
			float64(proccesInfo.AvailablePerfomance),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Capacity,
			prometheus.GaugeValue,
			float64(proccesInfo.Capacity),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryExcessTime,
			prometheus.GaugeValue,
			float64(proccesInfo.MemoryExcessTime),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SelectionSize,
			prometheus.GaugeValue,
			float64(proccesInfo.SelectionSize),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvgBackCallTime,
			prometheus.GaugeValue,
			float64(proccesInfo.AvgBackCallTime),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvgCallTime,
			prometheus.GaugeValue,
			float64(proccesInfo.AvgCallTime),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvgDbCallTime,
			prometheus.GaugeValue,
			float64(proccesInfo.AvgDbCallTime),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvgLockCallTime,
			prometheus.GaugeValue,
			float64(proccesInfo.AvgLockCallTime),
			proccesLabelsVal...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AvgServerCallTime,
			prometheus.GaugeValue,
			float64(proccesInfo.AvgServerCallTime),
			proccesLabelsVal...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Enable,
			prometheus.GaugeValue,
			func(fl bool) float64 {
				if fl {
					return 1.0
				} else {
					return 0.0
				}
			}(proccesInfo.Enable),
			proccesLabelsVal...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Running,
			prometheus.GaugeValue,
			func(fl bool) float64 {
				if fl {
					return 1.0
				} else {
					return 0.0
				}
			}(proccesInfo.Running),
			proccesLabelsVal...,
		)
		rpHostsCount++
	})

	ch <- prometheus.MustNewConstMetric(
		c.rpHosts,
		prometheus.GaugeValue,
		float64(rpHostsCount),
		clusterInfo.Name)

	c.wg.Done()
}
