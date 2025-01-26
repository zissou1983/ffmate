package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/yosev/debugo"
)

var gauges = make(map[string]prometheus.Gauge)
var gaugesVec = make(map[string]*prometheus.GaugeVec)
var dummyGauge = prometheus.NewGauge(prometheus.GaugeOpts{Name: "dummy", Help: "Dummy gauge"})
var dummyGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "dummy", Help: "Dummy gauge"}, []string{"dummy"})
var debug = debugo.New("prometheus:register")

type Metrics struct {
	Registry *prometheus.Registry
	Logger   *logrus.Logger
}

func (m *Metrics) Init() {
	m.Registry = prometheus.NewRegistry()
}

func (m *Metrics) RegisterGauge(name string, gauge prometheus.Gauge) {
	gauges[name] = gauge
	m.Registry.MustRegister(gauge)
	debug.Debugf("registered prometheus gauge '%s'", name)
}
func (m *Metrics) RegisterGaugeVec(name string, gauge *prometheus.GaugeVec) {
	gaugesVec[name] = gauge
	m.Registry.MustRegister(gauge)
	debug.Debugf("registered prometheus gaugeVec '%s'", name)
}

func (m *Metrics) Gauge(name string) prometheus.Gauge {
	if gauge, ok := gauges[name]; ok {
		return gauge
	}
	m.Logger.Warnf("prometheus gauge '%s' not found, returned dummy", name)
	return dummyGauge // return dummy gauge if gauge not found so .Inc() will never fail
}
func (m *Metrics) GaugeVec(name string) *prometheus.GaugeVec {
	if gauge, ok := gaugesVec[name]; ok {
		return gauge
	}
	m.Logger.Warnf("prometheus gaugeVec '%s' not found, returned dummy", name)
	return dummyGaugeVec // return dummy gauge if gauge not found so .Inc() will never fail
}

func (m *Metrics) Gauges() map[string]prometheus.Gauge {
	return gauges
}
func (m *Metrics) GaugesVec() map[string]*prometheus.GaugeVec {
	return gaugesVec
}
