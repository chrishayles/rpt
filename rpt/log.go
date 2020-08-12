package rpt

import (
	"fmt"
	"sync"
	"time"
)

// LOGGER

type Logger struct {
	LogOutputs    []Output
	MetricOutputs []Output
}

func (l *Logger) WriteLog(log *Log) {
	for _, o := range l.LogOutputs {
		o.WriteLog(log)
	}
}

func (l *Logger) WriteMetric(mc *MetricCollection) {
	for _, o := range l.MetricOutputs {
		o.WriteMetric(mc)
	}
}

func (l *Logger) AddLogOutput(o Output) {
	l.LogOutputs = append(l.LogOutputs, o)
}

func (l *Logger) AddMetricOutput(o Output) {
	l.MetricOutputs = append(l.MetricOutputs, o)
}

// OUTPUT

type Output interface {
	WriteLog(*Log)
	WriteMetric(*MetricCollection)
	Connect() error
	GetDescription() string
	resetMetrics()
	resetLogs()
	pullMetrics() *map[string]interface{}
	pullLogs() *map[string]interface{}
}

// LOG

type Log struct {
	Level       string // Debug, Info, Warn, Error
	Events      []*LogEvent
	Description string
}

type LogEvent struct {
	Level       string
	Time        time.Time
	Description string
}

func (l *Log) Debugf(format string, v ...interface{}) {
	if l.Level == "DEBUG" {
		l.logf("DEBUG", format, v...)
	}
}

func (l *Log) Infof(format string, v ...interface{}) {
	if l.Level == "DEBUG" || l.Level == "INFO" {
		l.logf("INFO", format, v...)
	}
}

func (l *Log) Warnf(format string, v ...interface{}) {
	if l.Level == "DEBUG" || l.Level == "INFO" || l.Level == "WARN" {
		l.logf("WARN", format, v...)
	}
}

func (l *Log) Errorf(format string, v ...interface{}) {
	l.logf("ERROR", format, v...)
}

func (l *Log) logf(level string, format string, v ...interface{}) {
	l.Events = append(l.Events, &LogEvent{
		Level:       level,
		Time:        time.Now(),
		Description: fmt.Sprintf(format, v...),
	})
}

func NewLog(level, description string) (*Log, error) {
	return nil, nil
}

// ELASTIC OUTPUT

type ElasticOutput struct {
	Description string
	ToProcess   chan *Log
}

func (e *ElasticOutput) WriteLog(l *Log) {}

func (e *ElasticOutput) WriteMetric(mc *MetricCollection) {}

func (e *ElasticOutput) Connect() error {
	return nil
}

func (e *ElasticOutput) GetDescription() string {
	return e.Description
}

func (e *ElasticOutput) resetLogs() {}

func (e *ElasticOutput) resetMetrics() {}

// FILE OUTPUT

type FileOutput struct {
	Description string
	ToProcess   chan *Log
	FileType    string //json, text, csv
	FilePath    string
	mu          sync.RWMutex
}

func (f *FileOutput) WriteLog(l *Log) {

	f.mu.Lock()

	//do stuff

	f.mu.Unlock()
}

func (f *FileOutput) WriteMetric(mc *MetricCollection) {

	f.mu.Lock()

	//do stuff

	f.mu.Unlock()
}

func (f *FileOutput) Connect() error {
	return nil
}

func (f *FileOutput) GetDescription() string {
	return f.Description
}

func (f *FileOutput) resetLogs() {}

func (f *FileOutput) resetMetrics() {}

// CONSOLE OUTPUT

type ConsoleOutput struct {
	Description string
	ToProcess   chan *Log
}

func (c *ConsoleOutput) WriteLog(l *Log) {

}

func (c *ConsoleOutput) WriteMetric(l *MetricCollection) {

}

func (c *ConsoleOutput) Connect() error {

	// No implementation needed.

	return nil
}

func (c *ConsoleOutput) GetDescription() string {

	return c.Description
}

func (c *ConsoleOutput) resetLogs() {}

func (c *ConsoleOutput) resetMetrics() {}

// PULL OUTPUT

type PullOutput struct {
	Description  string
	CacheMetrics []*MetricCollection
	CacheLogs    []*Log
}

func (p *PullOutput) WriteLog(l *Log) {
	p.CacheLogs = append(p.CacheLogs, l)
}

func (p *PullOutput) WriteMetric(mc *MetricCollection) {
	p.CacheMetrics = append(p.CacheMetrics, mc)
}

func (p *PullOutput) Connect() error {

	// No implementation needed.

	return nil
}

func (p *PullOutput) GetDescription() string {

	return p.Description
}

func (p *PullOutput) resetMetrics() {
	p.CacheMetrics = []*MetricCollection{}
}

func (p *PullOutput) resetLogs() {
	p.CacheLogs = []*Log{}
}

func (p *PullOutput) pullMetrics() *map[string]interface{} {

	output := &map[string]interface{}{
		"Metrics": p.CacheMetrics,
	}

	p.resetMetrics()

	return output
}

func (p *PullOutput) pullLogs() *map[string]interface{} {

	output := &map[string]interface{}{
		"Logs": p.CacheLogs,
	}

	p.resetLogs()

	return output
}

func NewPullOutput() *PullOutput {
	pull := &PullOutput{
		Description:  "pull_output",
		CacheLogs:    []*Log{},
		CacheMetrics: []*MetricCollection{},
	}

	return pull
}

// METRIC

type Metric struct {
	Label     string
	Value     interface{}
	Timestamp time.Time
}

type MetricCollection struct {
	Metrics []*Metric
}

func (mc *MetricCollection) AddMetric(m *Metric) {
	mc.Metrics = append(mc.Metrics, m)
}

func NewMetricCollection() (*MetricCollection, error) {

	m := []*Metric{}

	return &MetricCollection{
		Metrics: m,
	}, nil
}
