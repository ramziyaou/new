package entity

import (
	"sync"
	"time"
)

type CryptoWallet struct {
	Username string
	Name     string
	Amount   int64
	sync.RWMutex
}

type StartStopCheck struct {
	Username string
	Name     string
	Stop     bool
	Start    bool
}

func NewStartStop(username, name string) *StartStopCheck {
	return &StartStopCheck{
		username,
		name,
		false,
		false,
	}
}

func NewWallet(name string) *CryptoWallet {
	return &CryptoWallet{
		"",
		name,
		0,
		sync.RWMutex{},
	}
}

func (c *CryptoWallet) Mine() {
	time.Sleep(10 * time.Second)
	c.Lock()
	c.Amount++
	c.Unlock()
}
