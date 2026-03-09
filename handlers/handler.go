package handlers

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/tcodes0/jail-mcp/internal"
)

type job struct {
	id       string
	cmd      string
	started  time.Time
	stdout   bytes.Buffer
	stderr   bytes.Buffer
	exitCode int
	done     bool
	err      string
	mu       sync.Mutex
}

type Handler struct {
	cfg  *internal.Config
	jobs map[string]*job
	mu   sync.RWMutex
}

func New(cfg *internal.Config) *Handler {
	h := &Handler{
		cfg:  cfg,
		jobs: make(map[string]*job),
	}
	go h.removeJobsOlderThan(time.Hour)
	return h
}

func (h *Handler) removeJobsOlderThan(deadline time.Duration) {
	for range time.Tick(5 * time.Minute) {
		h.mu.Lock()
		for id, j := range h.jobs {
			j.mu.Lock()
			done, started := j.done, j.started
			j.mu.Unlock()
			if done && time.Since(started) > deadline {
				delete(h.jobs, id)
			}
		}
		h.mu.Unlock()
	}
}

// newJobID must be called with h.mu write lock held.
func (h *Handler) newJobID() string {
	id := fmt.Sprintf("%04d", rand.IntN(10000))
	if _, exists := h.jobs[id]; exists {
		return h.newJobID()
	}
	return id
}
