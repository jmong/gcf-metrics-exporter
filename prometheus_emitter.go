package metricsexporter
/**
 * Sends metrics to Prometheus PushGateway.
 * 
 * @usage
 * pusher := NewPrometheusPush(url, job)
 * pusher.Gatherer(myMetric)
 * pusher.Emit()
 **/

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/push"
)

const (
    PROM_PUSHGW_URL = "http://prometheus.monitoring:9091"
    PROM_PUSHGW_JOB = "pushgateway"
)

/* Prometheus Push object.
 */
type PrometheusPush struct {
    pusher  *push.Pusher
}

/* NewPrometheusPush creates a new Pusher object configured to 
 * send metrics to the Prometheus PushGateway located at the 
 * provided url and with the provided job name. 
 */
func NewPrometheusPush(url, job string) *PrometheusPush {
    return &PrometheusPush{
        pusher: push.New(url, job),
    }
}

/* Gatherer is just a facade around push.Gatherer().
 * See https://godoc.org/github.com/prometheus/client_golang/prometheus/push#Pusher.Gatherer
 */
func (p *PrometheusPush) Gatherer(g prometheus.Gatherer) {
    p.pusher.Gatherer(g)
}

/* Collector is just a facade around push.Collector().
 * See https://godoc.org/github.com/prometheus/client_golang/prometheus/push#Pusher.Collector
 */
func (p *PrometheusPush) Collector(c prometheus.Collector) {
    p.pusher.Collector(c)
}

/* Emit sends the metric to the Prometheus PushGateway.
 * It is just a wrapper around push.Push().
 * See https://godoc.org/github.com/prometheus/client_golang/prometheus/push#Pusher.Push
 */
func (p *PrometheusPush) Emit() error {
    return p.pusher.Push()
}
