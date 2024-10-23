package redis_rate_limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ()

func NewSortedSetCounterStrategy(client *redis.Client, now func() time.Time) Strategy {
	return &sortedSetCounter{
		client: client,
		now:    now,
	}
}

type sortedSetCounter struct {
	client *redis.Client
	now    func() time.Time
}

func (s *sortedSetCounter) Run(ctx context.Context, r *Request) (*Result, error) {
	now := s.now()

	item := uuid.New()
	minium := now.Add(-r.Duration)
	p := s.client.Pipeline()

	removeByScore := p.ZRemRangeByScore(ctx, r.Key, "0", strconv.FormatInt(minium.UnixMilli(), 10))

	add := p.ZAdd(ctx, r.Key, &redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: item.String(),
	})

	count := p.ZCount(ctx, r.Key, "-inf", "+inf")

	if _, err := p.Exec(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to execute sorted set pipeline for key: %v", r.Key)
	}
	if err := removeByScore.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to remove items from key %v", r.Key)
	}
	if err := add.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to add item to key %v", r.Key)
	}
	totalRequests, err := count.Result()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to count items for key %v", r.Key)
	}

	expiresAt := now.Add(r.Duration)
	requests := uint64(totalRequests)
	if requests > r.Limit {
		return &Result{
			State:         Deny,
			TotalRequests: requests,
			ExpiresAt:     expiresAt,
		}, nil
	}

	return &Result{
		State:         Allow,
		TotalRequests: requests,
		ExpiresAt:     expiresAt,
	}, nil
}
