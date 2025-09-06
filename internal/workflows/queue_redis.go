// Package workflows provides Redis-based adapters for workflows
package workflows

import (
	"encoding/json"
	"fmt"

	"github.com/archesai/archesai/internal/redis"
)

// QueueAdapter handles job queueing in Redis
type QueueAdapter struct {
	client *redis.Client
}

// NewQueueAdapter creates a new Redis queue adapter
func NewQueueAdapter(client *redis.Client) *QueueAdapter {
	return &QueueAdapter{
		client: client,
	}
}

// EnqueueRun adds a run to the processing queue
func (q *QueueAdapter) EnqueueRun(run *Run) error {
	queueName := fmt.Sprintf("workflows:runs:%s", run.Status)
	data, err := json.Marshal(run)
	if err != nil {
		return fmt.Errorf("failed to marshal run: %w", err)
	}
	return q.client.Queue.Push(queueName, data)
}

// DequeueRun retrieves and removes a run from the queue
func (q *QueueAdapter) DequeueRun(status string, timeout int) (*Run, error) {
	queueName := fmt.Sprintf("workflows:runs:%s", status)
	var run Run
	err := q.client.Queue.PopBlocking(queueName, timeout, &run)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

// PeekQueue checks the next item without removing it
func (q *QueueAdapter) PeekQueue(status string) (*Run, error) {
	queueName := fmt.Sprintf("workflows:runs:%s", status)
	var run Run
	err := q.client.Queue.Peek(queueName, &run)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

// GetQueueLength returns the number of items in a queue
func (q *QueueAdapter) GetQueueLength(status string) (int64, error) {
	queueName := fmt.Sprintf("workflows:runs:%s", status)
	return q.client.Queue.Length(queueName)
}

// MoveRun moves a run from one status queue to another
func (q *QueueAdapter) MoveRun(run *Run, fromStatus, toStatus string) error {
	// Remove from old queue
	fromQueue := fmt.Sprintf("workflows:runs:%s", fromStatus)
	if err := q.client.Queue.Remove(fromQueue, run); err != nil {
		return fmt.Errorf("failed to remove from %s queue: %w", fromStatus, err)
	}

	// Add to new queue with updated status
	run.Status = RunStatus(toStatus)
	return q.EnqueueRun(run)
}
