package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Monitoring struct {
	mu          sync.Mutex
	queryStats  map[string]*QueryStats
	alerts      []Alert
	alertConfig AlertConfig
}

type QueryStats struct {
	Count       int
	TotalTime   time.Duration
	ErrorCount  int
	LastError   string
}

type Alert struct {
	Timestamp time.Time
	Message   string
}

type AlertConfig struct {
	ErrorThreshold int
	TimeThreshold  time.Duration
}

func NewMonitoring() *Monitoring {
	return &Monitoring{
		queryStats: make(map[string]*QueryStats),
		alerts:     []Alert{},
		alertConfig: AlertConfig{
			ErrorThreshold: 5,
			TimeThreshold:  2 * time.Second,
		},
	}
}

func (m *Monitoring) RecordQuery(query string, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats, exists := m.queryStats[query]
	if !exists {
		stats = &QueryStats{}
		m.queryStats[query] = stats
	}

	stats.Count++
	stats.TotalTime += duration
	if err != nil {
		stats.ErrorCount++
		stats.LastError = err.Error()
	}

	if stats.ErrorCount >= m.alertConfig.ErrorThreshold || duration >= m.alertConfig.TimeThreshold {
		m.alerts = append(m.alerts, Alert{
			Timestamp: time.Now(),
			Message:   "Alert: Query " + query + " has high error rate or long duration",
		})
	}
}

func (m *Monitoring) GetQueryStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.queryStats)
}

func (m *Monitoring) GetAlerts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.alerts)
}

func (m *Monitoring) SetupRoutes(router *httprouter.Router) {
	router.GET("/monitoring/query-stats", m.GetQueryStats)
	router.GET("/monitoring/alerts", m.GetAlerts)
}
