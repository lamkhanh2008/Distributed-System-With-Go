package redis_rate_limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var (
	_ Strategy = &counterStrategy{}
)

const (
	keyWithoutExpire = -1
)

type counterStrategy struct {
	client *redis.Client
	now    func() time.Time
}

func NewCounterStrategy(client *redis.Client, now func() time.Time) *counterStrategy {
	return &counterStrategy{
		client: client,
		now:    now,
	}
}

func (c *counterStrategy) Run(ctx context.Context, r *Request) (*Result, error) {
	p := c.client.Pipeline()
	incrResult := p.Incr(ctx, r.Key)
	ttlResult := p.TTL(ctx, r.Key)

	if _, err := p.Exec(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to execute increment to key %v", r.Key)
	}

	totalRequests, err := incrResult.Result()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to increment key %v", r.Key)
	}
	var ttlDuration time.Duration

	if d, err := ttlResult.Result(); err != nil || d == keyWithoutExpire {
		ttlDuration = r.Duration
		if err := c.client.Expire(ctx, r.Key, r.Duration).Err(); err != nil {
			return nil, errors.Wrapf(err, "failed to set an expiration to key %v", r.Key)
		}
	} else {
		ttlDuration = d
	}

	expireAt := c.now().Add(ttlDuration)
	requests := uint64(totalRequests)

	if requests > r.Limit {
		return &Result{
			State:         Deny,
			TotalRequests: requests,
			ExpiresAt:     expireAt,
		}, nil
	}

	return &Result{
		State:         Allow,
		TotalRequests: requests,
		ExpiresAt:     expireAt,
	}, nil
}
