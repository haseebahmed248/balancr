// server pool management
package pool

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Backend struct {
	URL   string
	Alive bool
	mu    sync.Mutex
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

// GetNext returns next server (round-robin)
func (p *ServerPool) GetNext() (string, error) {
	p.mu.RLock()
	counter := 0
	defer p.mu.RUnlock()
	for {
		if counter == len(p.Backend) {
			return "", fmt.Errorf("No backend is alive")
		}
		backend := p.Backend[p.current]

		if backend.Alive != false {
			p.current = (p.current + 1) % len(p.Backend)
			return backend.URL, nil
		} else {
			p.current = (p.current + 1) % len(p.Backend)
		}
		counter++
	}
}
