// server pool management
package pool

import "sync"

type ServerPool struct {
	Backend []string
	current int
	mu      sync.Mutex
}

func GetServerPool(backend []string) *ServerPool {
	return &ServerPool{
		Backend: backend,
		current: 0,
	}
}

// GetNext returns next server (round-robin)
func (p *ServerPool) GetNext() string {
	p.mu.Lock()

	defer p.mu.Unlock()
	backend := p.Backend[p.current]
	if p.current == (len(p.Backend) - 1) {
		p.current = 0
	} else {
		p.current += 1
	}
	return backend
}
