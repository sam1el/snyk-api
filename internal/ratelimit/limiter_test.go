package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := Config{
		BurstSize:      5,
		Period:         200 * time.Millisecond,
		QueueSize:      100,
		MaxRetries:     3,
		RetryBaseDelay: 50 * time.Millisecond,
		RetryMaxDelay:  5 * time.Second,
	}

	limiter := New(cfg)
	require.NotNil(t, limiter)
	assert.Equal(t, cfg, limiter.config)
}

func TestLimiter_EnqueueAndExecute(t *testing.T) {
	cfg := Config{
		BurstSize: 2,
		Period:    100 * time.Millisecond,
		QueueSize: 10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limiter := New(cfg)
	limiter.Start(ctx, 1)
	defer limiter.Stop()

	var executed int
	var mu sync.Mutex

	// Enqueue multiple requests
	const numRequests = 5
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		req := &Request{
			ID:  "test-request",
			Ctx: ctx,
			Execute: func(ctx context.Context) error {
				mu.Lock()
				executed++
				mu.Unlock()
				return nil
			},
			Result: make(chan error, 1),
		}

		go func(r *Request) {
			defer wg.Done()
			limiter.Enqueue(r)
			<-r.Result
		}(req)
	}

	wg.Wait()

	mu.Lock()
	assert.Equal(t, numRequests, executed)
	mu.Unlock()
}

func TestLimiter_RateLimiting(t *testing.T) {
	cfg := Config{
		BurstSize: 2,
		Period:    200 * time.Millisecond, // 5 requests per second
		QueueSize: 10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limiter := New(cfg)
	limiter.Start(ctx, 1)
	defer limiter.Stop()

	startTime := time.Now()
	const numRequests = 4

	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		req := &Request{
			ID:  "rate-test",
			Ctx: ctx,
			Execute: func(ctx context.Context) error {
				return nil
			},
			Result: make(chan error, 1),
		}

		go func(r *Request) {
			defer wg.Done()
			limiter.Enqueue(r)
			<-r.Result
		}(req)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// First 2 requests should be immediate (burst)
	// Next 2 requests should be rate limited
	// Expect at least 200ms delay after burst
	assert.GreaterOrEqual(t, duration.Milliseconds(), int64(200))
}

func TestLimiter_ContextCancellation(t *testing.T) {
	cfg := DefaultConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter := New(cfg)
	limiter.Start(ctx, 1)

	reqCtx, reqCancel := context.WithCancel(context.Background())
	reqCancel() // Cancel immediately

	req := &Request{
		ID:  "cancelled-request",
		Ctx: reqCtx,
		Execute: func(ctx context.Context) error {
			t.Error("should not execute")
			return nil
		},
		Result: make(chan error, 1),
	}

	limiter.Enqueue(req)
	err := <-req.Result

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	limiter.Stop()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, 10, cfg.BurstSize)
	assert.Equal(t, 500*time.Millisecond, cfg.Period)
	assert.Equal(t, 1000, cfg.QueueSize)
	assert.Equal(t, 5, cfg.MaxRetries)
	assert.Equal(t, 100*time.Millisecond, cfg.RetryBaseDelay)
	assert.Equal(t, 30*time.Second, cfg.RetryMaxDelay)
}
