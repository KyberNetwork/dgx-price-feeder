package blockchain

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-data/common/blockchain"
	"github.com/KyberNetwork/reserve-data/common/blockchain/nonce"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	PRICING_OP string = "pricingOP"
)

type DGXReserve struct {
	*blockchain.BaseBlockchain
	reserve     *blockchain.Contract
	reserveAddr ethereum.Address
}

func (self *DGXReserve) GetAddresses() map[string]ethereum.Address {
	addrs := self.OperatorAddresses()
	addrs["dgx_reserve"] = self.reserveAddr
	return addrs
}

func (self *DGXReserve) RegisterPricingOperator(signer blockchain.Signer, nonceCorpus blockchain.NonceCorpus) {
	log.Printf("reserve pricing address: %s", signer.GetAddress().Hex())
	self.RegisterOperator(PRICING_OP, blockchain.NewOperator(signer, nonceCorpus))
}

//====================== Write calls ===============================

func (self *DGXReserve) SetPriceFeed(gasPrice *big.Int, blockNumber *big.Int, nonce *big.Int, ask1KDigix *big.Int, bid1KDigix *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	opts, err := self.GetTxOpts(PRICING_OP, nil, gasPrice, nil)
	if err != nil {
		return nil, err
	} else {
		timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		tx, err := self.BuildTx(timeout, opts, self.reserve, "setPriceFeed", blockNumber, nonce, ask1KDigix, bid1KDigix, v, r, s)
		if err != nil {
			return nil, err
		} else {
			return self.SignAndBroadcast(tx, PRICING_OP)
		}
	}
}

func (self *DGXReserve) Rebroadcast(tx *types.Transaction) (*types.Transaction, error) {
	return self.SignAndBroadcast(tx, PRICING_OP)
}

func NewDGXReserve(
	base *blockchain.BaseBlockchain,
	reserveAddr ethereum.Address,
	keystorePath string,
	passphrase string) *DGXReserve {

	log.Printf("reserve address: %s", reserveAddr.Hex())
	reserve := blockchain.NewContract(
		reserveAddr,
		"/go/src/github.com/KyberNetwork/dgx-price-feeder/blockchain/reserve.abi",
	)

	bc := &DGXReserve{
		BaseBlockchain: base,
		reserve:        reserve,
		reserveAddr:    reserveAddr,
	}

	signer := blockchain.NewEthereumSigner(keystorePath, passphrase)
	nonce := nonce.NewTimeWindow(signer.GetAddress())
	bc.RegisterPricingOperator(signer, nonce)

	return bc
}
