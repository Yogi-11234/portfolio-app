// Package ratelimit provides a shared RateLimiter interface and an in-memory
// sliding-window implementation. Keeping this in its own package breaks the
// import cycle between the middleware and service packages.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter is the contract for rate-limiting by an arbitrary string key.
type Limiter interface {
	// Allow reports whether a request identified by key is permitted given the
	// provided limit and time window. When the request is denied it also returns
	// how long the caller should wait before retrying.
	Allow(key string, limit int, window time.Duration) (allowed bool, retryAfter time.Duration)
}

// slidingWindowEntry stores the timestamps of accepted requests.
type slidingWindowEntry struct {
	timestamps []time.Time
}

type memoryLimiter struct {
	mu      sync.Mutex
	entries map[string]*slidingWindowEntry
}

// NewMemoryLimiter returns an in-memory sliding-window Limiter.
// It starts a background goroutine that purges stale entries every 5 minutes.
func NewMemoryLimiter() Limiter {
	l := &memoryLimiter{
		entries: make(map[string]*slidingWindowEntry),
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			l.mu.Lock()
			now := time.Now()
			for k, entry := range l.entries {
				if len(entry.timestamps) == 0 ||
					now.Sub(entry.timestamps[len(entry.timestamps)-1]) > time.Hour {
					delete(l.entries, k)
				}
			}
			l.mu.Unlock()
		}
	}()

	return l
}

func (m *memoryLimiter) Allow(key string, limit int, window time.Duration) (bool, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-window)

	entry, exists := m.entries[key]
	if !exists {
		m.entries[key] = &slidingWindowEntry{timestamps: []time.Time{now}}
		return true, 0
	}

	// Evict timestamps outside the window.
	valid := entry.timestamps[:0]
	for _, t := range entry.timestamps {
		if t.After(threshold) {
			valid = append(valid, t)
		}
	}
	entry.timestamps = valid

	if len(entry.timestamps) >= limit {
		retryAfter := entry.timestamps[0].Add(window).Sub(now)
		return false, retryAfter
	}

	entry.timestamps = append(entry.timestamps, now)
	return true, 0
}

// FormatRetryAfter formats a retryAfter duration as an integer seconds string
// suitable for the Retry-After HTTP header.
func FormatRetryAfter(d time.Duration) string {
	return fmt.Sprintf("%d", int(d.Seconds())+1)
}
