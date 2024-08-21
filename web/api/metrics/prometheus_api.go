package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

// @Author: morris
type PrometheusApi struct {
}

type KubeimoocCollector struct {
	clusterCpu prometheus.Gauge
	clusterMem prometheus.Gauge
}

func (k KubeimoocCollector) Describe(descs chan<- *prometheus.Desc) {
	k.clusterCpu.Describe(descs)
	k.clusterMem.Describe(descs)
}

func (k KubeimoocCollector) Collect(metrics chan<- prometheus.Metric) {
	usageArr := metricsService.GetClusterUsage()
	for _, item := range usageArr {
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

func newCollector() *KubeimoocCollector {
	return &KubeimoocCollector{
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

func (PrometheusApi) GetMetrics(ctx *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}
