package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

type PrometheusApi struct {
}

type KubemanageCollector struct {
	clusterCpu prometheus.Gauge
	clusterMem prometheus.Gauge
}

func (k KubemanageCollector) Describe(desc chan<- *prometheus.Desc) {
	k.clusterCpu.Describe(desc)
	k.clusterMem.Describe(desc)
}

// Collect 数据采集
func (k KubemanageCollector) Collect(metrics chan<- prometheus.Metric) {
	usageArray := metricsService.GetClusterUsage()
	for _, item := range usageArray {
		switch item.Label {
		case "cluster_cpu":
			newValue, _ := strconv.ParseFloat(item.Value, 64)
			k.clusterCpu.Set(newValue)
			k.clusterCpu.Collect(metrics)
		case "cluster_mem":
			newValue, _ := strconv.ParseFloat(item.Value, 64)
			k.clusterMem.Set(newValue)
			k.clusterMem.Collect(metrics)
		}

	}
}

func newCollector() *KubemanageCollector {
	return &KubemanageCollector{
		clusterCpu: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cluster_cpu",
				Help: "collector cluster cpu info",
			}),
		clusterMem: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cluster_mem",
				Help: "collector cluster memory info",
			}),
	}
}

func init() {
	prometheus.MustRegister(newCollector())
}

func (p *PrometheusApi) GetMetrics(ctx *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}
