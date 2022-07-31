package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	end   time.Time
}

func NewExecutionTimer() *ExecutionTimer {
	return &ExecutionTimer{
		histo: functionLatency,
		start: time.Now(),
		end:   time.Now(),
	}
}

//当前调用ObserveTotal时间减去start的时间放到直方图里面
func (e *ExecutionTimer) ObserveTotal() {
	e.histo.WithLabelValues("total").Observe(time.Now().Sub(e.start).Seconds())
}

var functionLatency = CreateExecutionTimeMetric("TimeExecutions", "Time spent.")

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

func RegisterMetrics() {
	if err := prometheus.Register(functionLatency); err != nil {
		panic(err)
	}
}
