package dgxpricing

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	NO_RETRY int = 6
	// INIT_GASPRICE int64 = 1000000000 // 1gwei
	INIT_GASPRICE int64 = 20000000000 // 20gwei
	GASPRICE_STEP int64 = 20000000000 // 20gwei
	NO_STEP       int   = 3
	// TX_WAIT_TIME  uint64 = 10 // 10 seconds
	TX_WAIT_TIME uint64 = 10 * 60 // 10 minutes
)

type PriceFeeder struct {
	runner  Runner
	reserve Reserve
	prices  PriceCorpus
}

func (self *PriceFeeder) Run() {
	self.runner.Start()
	self.feedPricePeriodically()
}

func (self *PriceFeeder) Stop() {
	self.runner.Stop()
}

func (self *PriceFeeder) TryFeedingPrice() (*types.Transaction, error) {
	blockno, nonce, ask, bid, v, r, s, err := self.prices.GetFeed()
	if err != nil {
		return nil, err
	}
	gasPrice := big.NewInt(INIT_GASPRICE)
	return self.reserve.SetPriceFeed(gasPrice, blockno, nonce, ask, bid, v, r, s)
}

func (self *PriceFeeder) MonitorAndRetry(tx *types.Transaction) error {
	startTime := time.Now()
	// this list should be sorted by gas price
	monitor := NewStatusMonitor(tx)
	for {
		// polling the tx status each 10s
		status, tx, err := monitor.GetStatus(self.reserve)
		if err != nil {
			log.Printf("Getting tx status failed: %s", err.Error())
		} else {
			switch status {
			case "pending":
				// it is still pending, if it is taking too long, replace it
				// with a new tx with higher nonce
				currentTime := time.Now()
				if uint64(currentTime.Sub(startTime)) > uint64(TX_WAIT_TIME*uint64(time.Second)) {
					newTx := types.NewTransaction(
						tx.Nonce(),
						*tx.To(),
						tx.Value(),
						tx.Gas(),
						big.NewInt(0).Add(tx.GasPrice(), big.NewInt(GASPRICE_STEP)),
						tx.Data(),
					)
					log.Printf("Replacing old tx with tx %s", newTx.Hash().Hex())
					newSignedTx, err := self.reserve.Rebroadcast(newTx)
					if err != nil {
						log.Printf("Replacing old tx failed, err(%s). Ignore, will try next time", err.Error())
					} else {
						monitor.PushTx(newSignedTx)
					}
				}
			case "lost":
				// retry
				_, err := self.reserve.Rebroadcast(tx)
				if err != nil {
					log.Printf("Rebroadcasting the tx failed, err(%s). Ignore, will try next time", err.Error())
				}
			case "mined":
				// the tx is successfully done
				log.Printf("Tx %s is mined. Finish monitoring.", tx.Hash().Hex())
				return nil
			case "failed":
				// we dont retry in this case, it will just fail
				log.Printf("Tx %s is failed. Finish monitoring.", tx.Hash().Hex())
				return nil
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (self *PriceFeeder) EnsureFeedPrice() {
	for i := 0; i < NO_RETRY; i++ {
		log.Printf("Try feeding price")
		tx, err := self.TryFeedingPrice()
		if err != nil {
			log.Printf("%d(th) Try failed: err(%s)", err.Error())
		} else {
			// monitor the status and increase the gas price by GASPRICE_STEP if needed
			// gas price will be increased only NO_STEP times

			// err will be returned only when we giveup increasing nonce
			err := self.MonitorAndRetry(tx)
			if err != nil {
				log.Printf("Gave up on setting the price feed after increased the nonce for %d times", NO_STEP)
			}
			return
		}
	}
}

func (self *PriceFeeder) feedPricePeriodically() {
	for {
		log.Printf("Going to feed the price to the contract")
		self.EnsureFeedPrice()
		log.Printf("Waiting for signal for the next interval...")
		<-self.runner.GetPricingTicker()
	}
}

func NewPriceFeeder(runner Runner, reserve Reserve, prices PriceCorpus) *PriceFeeder {
	return &PriceFeeder{
		runner:  runner,
		reserve: reserve,
		prices:  prices,
	}
}
