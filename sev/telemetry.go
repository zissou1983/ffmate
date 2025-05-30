package sev

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promDto "github.com/prometheus/client_model/go"
	"github.com/welovemedia/ffmate/internal/config"
)

type Stats struct {
	AppName        string `json:"appName"`
	AppVersion     string `json:"appVersion"`
	ClientId       string `json:"clientId"`
	SessionId      string `json:"sessionId"`
	IsShuttingDown bool   `json:"isShuttingDown"`

	RuntimeDuration int64 `json:"runtimeDuration"`

	Os   string `json:"os"`
	Arch string `json:"arch"`

	Metrics map[string]float64

	Stats  map[string]interface{}
	Config map[string]interface{}
}

var debugTelemetry = debug.Extend("telemetry")

func (s *Sev) SendTelemetry(targetUrl string, statistics map[string]interface{}, conf map[string]interface{}) {
	stats := Stats{
		AppName:        config.Config().AppName,
		AppVersion:     config.Config().AppVersion,
		ClientId:       s.Client().Uuid,
		SessionId:      s.Session(),
		IsShuttingDown: s.isShuttingDown,

		RuntimeDuration: time.Since(s.AppStartTime()).Milliseconds(),

		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,

		Metrics: make(map[string]float64),

		Stats:  statistics,
		Config: conf,
	}

	for name, gauge := range s.Metrics().Gauges() {
		g := &promDto.Metric{}
		gauge.Write(g)
		stats.Metrics[name] = g.Gauge.GetValue()
	}
	for name, gaugeVec := range s.Metrics().GaugesVec() {
		metricChan := make(chan prometheus.Metric, 1)

		go func() {
			gaugeVec.Collect(metricChan)
			close(metricChan)
		}()

		for metric := range metricChan {
			promMetric := &promDto.Metric{}
			if err := metric.Write(promMetric); err != nil {
				fmt.Printf("Error writing metric: %v\n", err)
				continue
			}

			labelValues := make([]string, len(promMetric.Label))
			for i, label := range promMetric.Label {
				labelValues[i] = fmt.Sprintf("%s=%s", *label.Name, *label.Value)
			}

			labeledName := fmt.Sprintf("%s{%s}", name, strings.Join(labelValues, ","))
			stats.Metrics[labeledName] = promMetric.Gauge.GetValue()
		}
	}

	b, err := json.Marshal(&stats)
	if err != nil {
		s.Logger().Warnf("failed to marshal telemetry data: %+v", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(b))
	if err != nil {
		s.Logger().Error("failed to create http request", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", config.Config().AppName+"/"+config.Config().AppVersion)

	_, err = client.Do(req)
	if err != nil {
		s.Logger().Warnf("failed to send telemetry data: %+v", err)
	} else {
		debugTelemetry.Debugf("sent telemetry data")
	}
}
