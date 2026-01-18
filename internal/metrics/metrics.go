package metrics

import (
	"fmt"
	"sync"
)

type Metrics struct {
	TotalRequests      int64
	RequestsPerBackend (map[string]int64) // backendUrl: count
	Errors             int64
	mu                 sync.Mutex
}

func (m *Metrics) IncrementRequest(backendURL string) {
	m.mu.Lock()
	if m.RequestsPerBackend == nil {
		m.RequestsPerBackend = make(map[string]int64)
	}
	defer m.mu.Unlock()
	if _, ok := m.RequestsPerBackend[backendURL]; ok {
		m.RequestsPerBackend[backendURL]++
		m.TotalRequests++
		return
	}
	m.RequestsPerBackend[backendURL] = 1
	m.TotalRequests++
}

func (m *Metrics) IncrementError() {
	m.mu.Lock()
	if m.RequestsPerBackend == nil {
		m.RequestsPerBackend = make(map[string]int64)
	}
	defer m.mu.Unlock()
	m.Errors++
	m.TotalRequests++
}

func (m *Metrics) GetMetrics() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := fmt.Sprintf("Total: %d | Errors: %d\n", m.TotalRequests, m.Errors)
	for backend, count := range m.RequestsPerBackend {
		data += fmt.Sprintf("%s: %d\n", backend, count)
	}
	return data
}
