package ratelimit

import (
	"errors"
	"github.com/go-kratos/aegis/ratelimit"
	"sync"
	"time"
)

// TokenBucketLimiter 基于令牌桶算法的限流器实现
type TokenBucketLimiter struct {
	rate       int        // 每秒产生的令牌数
	capacity   int        // 桶的容量
	tokens     int        // 当前令牌数
	lastUpdate time.Time  // 上次更新时间
	mu         sync.Mutex // 互斥锁保证并发安全
}

// DoneFunc is done function.

// DoneInfo is done info.

// NewTokenBucketLimiter 创建令牌桶限流器
// rate: 每秒允许的请求数
// capacity: 桶容量（突发流量允许的最大请求数）
func NewTokenBucketLimiter(rate, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastUpdate: time.Now(),
	}
}

// Allow 实现 Limiter 接口
func (l *TokenBucketLimiter) Allow() (ratelimit.DoneFunc, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 计算当前应补充的令牌数（按时间流逝比例）
	now := time.Now()
	elapsed := now.Sub(l.lastUpdate)
	tokensToAdd := int(elapsed.Seconds() * float64(l.rate))
	if tokensToAdd > 0 {
		l.tokens = min(l.tokens+tokensToAdd, l.capacity)
		l.lastUpdate = now
	}

	// 令牌不足时拒绝请求
	if l.tokens <= 0 {
		return nil, ErrLimitExceed
	}

	// 消耗令牌
	l.tokens--
	return func(ratelimit.DoneInfo) {}, nil
}

// DoneFunc 空实现（令牌桶算法无需回收令牌）
func (l *TokenBucketLimiter) Done() {}

// 工具函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 错误定义
var ErrLimitExceed = errors.New("rate limit exceeded")
