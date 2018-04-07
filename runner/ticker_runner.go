package runner

import (
	"time"
)

type TickerRunner struct {
	duration time.Duration
	clock    *time.Ticker
	signal   chan bool
}

func (self *TickerRunner) GetPricingTicker() <-chan time.Time {
	if self.clock == nil {
		<-self.signal
	}
	return self.clock.C
}

func (self *TickerRunner) Start() error {
	self.clock = time.NewTicker(self.duration)
	self.signal <- true
	return nil
}

func (self *TickerRunner) Stop() error {
	self.clock.Stop()
	return nil
}

func NewTickerRunner(duration time.Duration) *TickerRunner {
	return &TickerRunner{
		duration,
		nil,
		make(chan bool, 1),
	}
}
