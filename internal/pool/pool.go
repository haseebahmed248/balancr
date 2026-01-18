// server pool management
package pool

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Backend struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
	Alive  bool
	Quota  int
	mu     sync.Mutex
}

type ServerPool struct {
	Backend []*Backend
	current int
	mu      sync.RWMutex
}

func GetServerPool(backend []*Backend) *ServerPool {
	return &ServerPool{
		Backend: backend,
		current: 0,
	}
}

// SetAlive
func SetAlive(backend *Backend) {
	backend.Alive = true
}

// IsAlive GET
func IsAlive(backend *Backend) bool {
	conn, err := net.DialTimeout("tcp", backend.URL, 2*time.Second)
	backend.mu.Lock()
	if err != nil {
		backend.Alive = false
		backend.mu.Unlock()
		return false
	}

	defer conn.Close()
	defer backend.mu.Unlock()
	if backend.Alive == false {
		SetAlive(backend)
	}
	return true
}

// Health Check (for Now using net.DialTImeout but will upgrade to /health with net/http)
func (p *ServerPool) HealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()
	for range ticker.C {
		for _, backend := range p.Backend {
			go IsAlive(backend)
		}
	}

}

func resetQuota(backends []*Backend) {
	for _, backend := range backends {
		backend.Quota = backend.Weight
	}
}

// GetNext returns next server (round-robin)
func (p *ServerPool) GetNext() (string, error) {
	p.mu.RLock()
	counter := 0
	var hasAlive bool
	defer p.mu.RUnlock()
	for {
		backend := p.Backend[p.current]
		if counter == len(p.Backend) {
			if hasAlive == false {
				return "", fmt.Errorf("No backend is alive")
			} else {
				resetQuota(p.Backend)
				counter = 0
				hasAlive = false
				continue // restart the loop
			}
		}
		if backend.Alive != false && backend.Quota > 0 {
			hasAlive = true
			p.current = (p.current + 1) % len(p.Backend)
			backend.Quota--
			return backend.URL, nil
		} else {
			p.current = (p.current + 1) % len(p.Backend)
			if backend.Alive {
				hasAlive = true
			}
		}
		counter++
	}
}
