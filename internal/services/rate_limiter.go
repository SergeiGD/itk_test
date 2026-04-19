package services

import (
	"sync"
	"time"

	"github.com/SergeiGD/itk_test/config"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type WalletRateLimiter interface {
	IsAllowed(wallet uuid.UUID) bool
}

type limiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newLimiter(l *rate.Limiter) *limiter {
	return &limiter{
		limiter:  l,
		lastSeen: time.Now(),
	}
}

type walletRateLimiter struct {
	mu   sync.Mutex
	mp   map[uuid.UUID]*limiter
	conf config.Config
}

func NewWalletRateLimiter(conf config.Config) WalletRateLimiter {
	l := &walletRateLimiter{
		mp:   make(map[uuid.UUID]*limiter, conf.Limiter.MaxLimit),
		mu:   sync.Mutex{},
		conf: conf,
	}
	go l.startCleaner()
	return l
}

func (rl *walletRateLimiter) IsAllowed(wallet uuid.UUID) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if l, ok := rl.mp[wallet]; ok {
		return l.limiter.Allow()
	}

	l := newLimiter(rate.NewLimiter(rate.Limit(rl.conf.Limiter.MaxLimit), rl.conf.Limiter.Burst))
	rl.mp[wallet] = l
	return l.limiter.Allow()
}

func (rl *walletRateLimiter) startCleaner() {
	for {
		time.Sleep(rl.conf.Limiter.CleanInterval)
		rl.mu.Lock()
		for wallet, l := range rl.mp {
			if time.Since(l.lastSeen) > 5*time.Minute {
				delete(rl.mp, wallet)
			}
		}
		rl.mu.Unlock()
	}
}
