package mutablelimiter

import "context"

type Limiter struct {
	adjustLimit chan int
	acquire     chan acquireRequest
	getLimit    chan struct{ cap, len int }
}

type acquireResponse struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type acquireRequest struct {
	ctx  context.Context
	resp chan<- acquireResponse
}
