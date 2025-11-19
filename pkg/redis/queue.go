package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Queue provides queue operations using Redis lists.
type Queue struct {
	client *redis.Client
}

// NewQueue creates a new queue instance.
func NewQueue(client *redis.Client) *Queue {
	return &Queue{
		client: client,
	}
}

// Push adds an item to the queue.
func (q *Queue) Push(queueName string, item any) error {
	ctx := context.Background()
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	return q.client.RPush(ctx, queueName, data).Err()
}

// Pop removes and returns an item from the queue.
func (q *Queue) Pop(queueName string, dest any) error {
	ctx := context.Background()
	data, err := q.client.LPop(ctx, queueName).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// PopBlocking removes and returns an item from the queue, blocking if empty.
func (q *Queue) PopBlocking(queueName string, timeout int, dest any) error {
	ctx := context.Background()
	result, err := q.client.BLPop(ctx, time.Duration(timeout)*time.Second, queueName).Result()
	if err != nil {
		return err
	}

	if len(result) < 2 {
		return fmt.Errorf("unexpected result from BLPop")
	}

	return json.Unmarshal([]byte(result[1]), dest)
}

// Peek returns the next item without removing it.
func (q *Queue) Peek(queueName string, dest any) error {
	ctx := context.Background()
	data, err := q.client.LIndex(ctx, queueName, 0).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Length returns the queue length.
func (q *Queue) Length(queueName string) (int64, error) {
	ctx := context.Background()
	return q.client.LLen(ctx, queueName).Result()
}

// Clear removes all items from the queue.
func (q *Queue) Clear(queueName string) error {
	ctx := context.Background()
	return q.client.Del(ctx, queueName).Err()
}

// Remove removes a specific item from the queue.
func (q *Queue) Remove(queueName string, item any) error {
	ctx := context.Background()
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	// Remove 1 occurrence of the item
	removed, err := q.client.LRem(ctx, queueName, 1, string(data)).Result()
	if err != nil {
		return err
	}
	if removed == 0 {
		return fmt.Errorf("item not found in queue")
	}
	return nil
}
