package pool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(5)
	defer func() { _ = pool.Close() }()

	var counter int32
	for i := 0; i < 10; i++ {
		err := pool.Submit(func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
		if err != nil {
			t.Errorf("failed to submit task: %v", err)
		}
	}

	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&counter) != 10 {
		t.Errorf("expected 10 tasks completed, got %d", counter)
	}
}

func TestPoolClose(t *testing.T) {
	pool := NewWorkerPool(2)
	err := pool.Close()
	if err != nil {
		t.Errorf("close failed: %v", err)
	}

	err = pool.Submit(func(ctx context.Context) error {
		return nil
	})
	if err == nil {
		t.Error("should reject tasks after close")
	}
}
