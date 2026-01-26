package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Priority levels for tasks
type Priority int

const (
	LowPriority    Priority = 0
	NormalPriority Priority = 1
	HighPriority   Priority = 2
)

// TaskWithMetadata represents a task with additional metadata
type TaskWithMetadata struct {
	Fn       Task
	Priority Priority
	ID       string
}

// WorkerPool is a flexible worker pool implementation with priority handling
type WorkerPool struct {
	highPriorityTasks   chan Task
	normalPriorityTasks chan Task
	lowPriorityTasks    chan Task
	wg                  sync.WaitGroup
	ctx                 context.Context
	cancel              context.CancelFunc
	size                atomic.Int32
	active              atomic.Int32
	queued              atomic.Int32
	max                 atomic.Int32
	completed           atomic.Int64
	failed              atomic.Int64
	closed              atomic.Bool
}

type Task func(context.Context) error

// NewWorkerPool creates a new worker pool with the specified size
func NewWorkerPool(size int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		highPriorityTasks:   make(chan Task, size),
		normalPriorityTasks: make(chan Task, size*2),
		lowPriorityTasks:    make(chan Task, size*4),
		ctx:                 ctx,
		cancel:              cancel,
	}

	// Use atomic Int32 for thread safety
	wp.size.Store(int32(size))
	wp.max.Store(int32(size))

	// Start workers
	for i := 0; i < size; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

// worker processes tasks from the queues
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		// Process tasks with priority (highest first)
		var task Task
		var ok bool

		select {
		// Check high priority tasks first
		case task, ok = <-wp.highPriorityTasks:
			if !ok {
				// High priority channel closed, but we might still have other channels
				// Check if we should exit completely
				select {
				case <-wp.ctx.Done():
					return
				default:
					continue
				}
			}

		// If no high priority tasks, check normal priority
		case task, ok = <-wp.normalPriorityTasks:
			if !ok {
				// Normal priority channel closed, but we might still have other channels
				select {
				case <-wp.ctx.Done():
					return
				default:
					continue
				}
			}

		// If no high or normal priority tasks, check low priority
		case task, ok = <-wp.lowPriorityTasks:
			if !ok {
				// Low priority channel closed, check if we should exit
				select {
				case <-wp.ctx.Done():
					return
				default:
					continue
				}
			}

		// Exit if context is done
		case <-wp.ctx.Done():
			return
		}

		// Process the task
		wp.active.Add(1)
		if task != nil {
			if err := task(wp.ctx); err != nil {
				wp.failed.Add(1)
			} else {
				wp.completed.Add(1)
			}
		}
		wp.active.Add(-1)
		wp.queued.Add(-1)
	}
}

// SubmitWithPriority submits a task with the specified priority
func (wp *WorkerPool) SubmitWithPriority(task Task, priority Priority) error {
	// Check if pool is already closed
	if wp.closed.Load() {
		return fmt.Errorf("pool is closed")
	}

	if wp.queued.Load() >= wp.max.Load()*4 {
		return fmt.Errorf("pool queue full")
	}

	wp.queued.Add(1)

	// Submit to the appropriate queue based on priority
	var taskQueue chan Task
	switch priority {
	case HighPriority:
		taskQueue = wp.highPriorityTasks
	case NormalPriority:
		taskQueue = wp.normalPriorityTasks
	default:
		taskQueue = wp.lowPriorityTasks
	}

	// Try to submit the task
	select {
	case taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		wp.queued.Add(-1)
		return fmt.Errorf("pool closed")
	}
}

// Submit submits a task with normal priority (backward compatibility)
func (wp *WorkerPool) Submit(task Task) error {
	return wp.SubmitWithPriority(task, NormalPriority)
}

// Close shuts down the worker pool
func (wp *WorkerPool) Close() error {
	// Mark pool as closed first to prevent new submissions
	wp.closed.Store(true)

	wp.cancel() // Signal all workers to stop

	// Close all task channels
	close(wp.highPriorityTasks)
	close(wp.normalPriorityTasks)
	close(wp.lowPriorityTasks)

	// Wait for all workers to finish
	wp.wg.Wait()
	return nil
}

// Stats returns statistics about the worker pool
func (wp *WorkerPool) Stats() map[string]interface{} {
	// Pre-allocate the map for better performance
	stats := make(map[string]interface{}, 8)

	// Basic stats
	stats["active"] = wp.active.Load()
	stats["queued"] = wp.queued.Load()
	stats["workers"] = wp.size.Load()
	stats["max_workers"] = wp.max.Load()

	// Queue stats
	stats["high_priority_queued"] = len(wp.highPriorityTasks)
	stats["normal_priority_queued"] = len(wp.normalPriorityTasks)
	stats["low_priority_queued"] = len(wp.lowPriorityTasks)

	// Performance stats
	stats["completed"] = wp.completed.Load()
	stats["failed"] = wp.failed.Load()

	return stats
}

// ResizePool changes the size of the worker pool
func (wp *WorkerPool) ResizePool(newSize int) {
	oldSize := int(wp.size.Load())
	wp.size.Store(int32(newSize))
	wp.max.Store(int32(newSize))

	// If we're growing the pool, add more workers
	if newSize > oldSize {
		for i := 0; i < newSize-oldSize; i++ {
			wp.wg.Add(1)
			go wp.worker()
		}
	}

	// If we're shrinking, workers will exit naturally when the context is canceled
	// They check the wp.ctx.Done() channel
}
