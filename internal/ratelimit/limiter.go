// Package ratelimit provides token bucket rate limiting with worker pools
// for managing concurrent API requests with configurable burst and period settings.
package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Config holds configuration for the rate limiter.
type Config struct {
	BurstSize      int           // Maximum concurrent requests
	Period         time.Duration // Time between allowing new requests
	QueueSize      int           // Buffer size for pending requests
	MaxRetries     int           // Maximum retry attempts
	RetryBaseDelay time.Duration // Initial retry delay
	RetryMaxDelay  time.Duration // Maximum retry delay
}

// DefaultConfig returns sensible defaults for rate limiting.
func DefaultConfig() Config {
	return Config{
		BurstSize:      10,
		Period:         500 * time.Millisecond,
		QueueSize:      1000,
		MaxRetries:     5,
		RetryBaseDelay: 100 * time.Millisecond,
		RetryMaxDelay:  30 * time.Second,
	}
}

// Limiter implements rate limiting using a token bucket algorithm.
type Limiter struct {
	limiter *rate.Limiter
	queue   chan *Request
	wg      sync.WaitGroup
	done    chan struct{}
	config  Config
}

// Request represents a queued request waiting for rate limit permission.
type Request struct {
	ID      string
	Execute func(context.Context) error
	Result  chan error
	Ctx     context.Context
}

// New creates a new rate limiter with the given configuration.
func New(cfg Config) *Limiter {
	// Calculate rate: requests per second
	// e.g., period=500ms means 2 requests per second
	requestsPerSecond := rate.Limit(float64(time.Second) / float64(cfg.Period))

	if cfg.QueueSize < 1 {
		cfg.QueueSize = 1000
	}

	return &Limiter{
		limiter: rate.NewLimiter(requestsPerSecond, cfg.BurstSize),
		queue:   make(chan *Request, cfg.QueueSize),
		done:    make(chan struct{}),
		config:  cfg,
	}
}

// Start begins processing the queue with the specified number of workers.
func (l *Limiter) Start(ctx context.Context, workerCount int) {
	if workerCount < 1 {
		workerCount = 1
	}

	for i := 0; i < workerCount; i++ {
		l.wg.Add(1)
		go l.worker(ctx)
	}
}

// worker processes requests from the queue with rate limiting.
func (l *Limiter) worker(ctx context.Context) {
	defer l.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-l.done:
			return
		case req := <-l.queue:
			// Wait for rate limiter permission
			// nolint:contextcheck // Intentionally using request context, not worker context
			if err := l.limiter.Wait(req.Ctx); err != nil {
				// Context cancelled
				req.Result <- err
				continue
			}

			// Execute the request
			// nolint:contextcheck // Intentionally using request context, not worker context
			err := req.Execute(req.Ctx)
			req.Result <- err
		}
	}
}

// Enqueue adds a request to the queue for rate-limited execution.
func (l *Limiter) Enqueue(req *Request) {
	l.queue <- req
}

// Stop gracefully stops the rate limiter.
func (l *Limiter) Stop() {
	close(l.done)
	l.wg.Wait()
}
