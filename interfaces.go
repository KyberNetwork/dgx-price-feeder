package dgxpricing

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Runner interface {
	GetPricingTicker() <-chan time.Time
	Start() error
	Stop() error
}

type PriceCorpus interface {
	GetFeed() (blockNumber *big.Int, nonce *big.Int, ask1KDigix *big.Int, bid1KDigix *big.Int, v uint8, r [32]byte, s [32]byte, err error)
}

type Reserve interface {
	SetPriceFeed(gasPrice *big.Int, blockNumber *big.Int, nonce *big.Int, ask1KDigix *big.Int, bid1KDigix *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error)
	TxStatus(common.Hash) (status string, blockno uint64, err error)
	Rebroadcast(tx *types.Transaction) (*types.Transaction, error)
}
