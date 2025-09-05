package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

// Queue provides queue operations using Redis lists
type Queue struct{}

// NewQueue creates a new queue instance
func NewQueue() *Queue {
	return &Queue{}
}

// Push adds an item to the queue
func (q *Queue) Push(queueName string, item interface{}) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	return client.RPush(ctx, queueName, data).Err()
}

// Pop removes and returns an item from the queue
func (q *Queue) Pop(queueName string, dest interface{}) error {
	data, err := client.LPop(ctx, queueName).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// PopBlocking removes and returns an item from the queue, blocking if empty
func (q *Queue) PopBlocking(queueName string, timeout int, dest interface{}) error {
	result, err := client.BLPop(ctx, time.Duration(timeout)*time.Second, queueName).Result()
	if err != nil {
		return err
	}

	if len(result) < 2 {
		return fmt.Errorf("unexpected result from BLPop")
	}

	return json.Unmarshal([]byte(result[1]), dest)
}

// Peek returns the next item without removing it
func (q *Queue) Peek(queueName string, dest interface{}) error {
	data, err := client.LIndex(ctx, queueName, 0).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Length returns the queue length
func (q *Queue) Length(queueName string) (int64, error) {
	return client.LLen(ctx, queueName).Result()
}

// Clear removes all items from the queue
func (q *Queue) Clear(queueName string) error {
	return client.Del(ctx, queueName).Err()
}
