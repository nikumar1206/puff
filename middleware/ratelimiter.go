package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nikumar1206/puff"
	"github.com/redis/go-redis/v9"
)

// RateLimiterType defines an interface for implementations of a
// rate limiter.
type RateLimiterType interface {
	// Accept specifies a function that takes in the client detail, the
	// client's limit, and the window for how long to store the request.
	// It should return a boolean on whether or not to allow the request.
	// Defaults to the RateLimiterInMemoryAccept.
	Accept(key string, limit int, timeWindow time.Duration) bool
}

// RateLimiterInMemory implements the RateLimiterType interface
// and stores client request information in memory.
type RateLimiterInMemory struct {
	requests map[string]int
}

func (i *RateLimiterInMemory) Accept(key string, limit int, timeWindow time.Duration) bool {
	if i.requests == nil {
		i.requests = make(map[string]int)
	}
	if limit == 0 { // no need to check
		return false
	}
	if i.requests[key] >= limit {
		return false
	}
	i.requests[key] += 1
	go func() {
		time.Sleep(timeWindow)
		i.requests[key] = i.requests[key] - 1
	}()
	return true
}

// RateLimiterInRedis implements a rate limiter in redis.
// The redis client must be set in order for the rate limiter
// to work.
type RateLimiterInRedis struct {
	RedisClient *redis.Client
}

func (r *RateLimiterInRedis) Accept(key string, limit int, timeWindow time.Duration) bool {
	if r.RedisClient == nil {
		slog.Error("redis client is nil, allowing request to continue.")
		return true
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	ll := r.RedisClient.LPush(ctx, key, time.Now().UnixNano())
	if ll == nil {
		slog.Error(fmt.Sprintf("failed to get length of list in redis with key %s", key))
		return true // allow anyway
	}
	if ll.Val() > int64(limit) {
		r.RedisClient.LPop(context.Background(), key)
		return false
	}
	go func() {
		time.Sleep(timeWindow)
		r.RedisClient.LPop(context.Background(), key)
	}()
	return true
}

// func (*r)

type RateLimiterConfig struct {
	// RateLimiter specifies a value that implements the RateLimiterType interface.
	RateLimiter RateLimiterType
	// RequestTimeWindow specifies how long until a request from a client
	// should be forgotten. Once a request from a client comes in, the rate
	// limiter looks for requests from the client that it has not yet forgotten.
	// If the number of requests that were found exceeds, the maximum number
	// allowed, it will return to the user a 429 Too Many Requests. Defaults to
	// one minute.
	RequestTimeWindow time.Duration
	// GetClientDetail should return a key that can be used to identify the client
	// making requests. Defaults to client's IP Address.
	GetClientDetail func(*puff.Context) string
	// GetClientRateLimit is a specifiable function that returns the rate limit
	// for a client based on the context of the request. If this function returns 0,
	// the request is blocked since no matter how many requests have come in from the
	// client, it exceeds it anyway.  Possible use cases include: authenticated users
	// getting more requests to your server than non-authenticated users. Defaults to 60.
	GetClientRateLimit func(*puff.Context) int
	// GetRateLimitedResponse is a specifiable function that specifies that response to give
	// if the client has been rate limited.
	GetRateLimitedResponse func(*puff.Context) puff.Response
}

// DefaultRateLimiterConfig provides the default configuration for the RateLimiter middleware.
var DefaultRateLimiterConfig RateLimiterConfig = RateLimiterConfig{
	RateLimiter:       new(RateLimiterInMemory),
	RequestTimeWindow: time.Minute,
	GetClientDetail: func(c *puff.Context) string {
		return c.ClientIP()
	},
	GetClientRateLimit: func(_ *puff.Context) int {
		return 60
	},
	GetRateLimitedResponse: func(_ *puff.Context) puff.Response {
		return puff.JSONResponse{
			StatusCode: 429,
			Content: map[string]string{
				"error": "You have been rate limited.",
			},
		}
	},
}

// createRateLimiterMiddleware creates a rate limiter middleware with the given configuration.
func createRateLimiterMiddleware(r RateLimiterConfig) puff.Middleware {
	return func(next puff.HandlerFunc) puff.HandlerFunc {
		return func(ctx *puff.Context) {
			if r.RateLimiter.Accept(r.GetClientDetail(ctx), r.GetClientRateLimit(ctx), r.RequestTimeWindow) {
				next(ctx)
				return
			}
			ctx.SendResponse(r.GetRateLimitedResponse(ctx))
		}
	}
}

// RateLimiter returns a RateLimiter middleware with the default config.
func RateLimiter() puff.Middleware {
	return createRateLimiterMiddleware(DefaultRateLimiterConfig)
}

// RateLimiterWithConfig returns a RateLimiter middleware with the specified configuration.
func RateLimiterWithConfig(r RateLimiterConfig) puff.Middleware {
	return createRateLimiterMiddleware(r)
}
