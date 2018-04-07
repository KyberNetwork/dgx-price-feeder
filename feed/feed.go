package feed

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	// "strings"
	"log"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const ENDPOINT string = "http://www.9gum3.com/feed"

type pricejson struct {
	Block   *big.Int         `json:"block_number"`
	Ask     *big.Int         `json:"ask_for_1000"`
	Bid     *big.Int         `json:"bid_for_1000"`
	Hash    ethereum.Hash    `json:"hash"`
	Nonce   *big.Int         `json:"nonce"`
	R       string           `json:"r"`
	S       string           `json:"s"`
	V       uint8            `json:"v"`
	Signer  ethereum.Address `json:"signer"`
	Message string           `json:"message"`
}

type Price struct {
	Block   *big.Int         `json:"block_number"`
	Ask     *big.Int         `json:"ask_for_1000"`
	Bid     *big.Int         `json:"bid_for_1000"`
	Hash    ethereum.Hash    `json:"hash"`
	Nonce   *big.Int         `json:"nonce"`
	R       [32]byte         `json:"r"`
	S       [32]byte         `json:"s"`
	V       uint8            `json:"v"`
	Signer  ethereum.Address `json:"signer"`
	Message string           `json:"message"`
}

type PriceFeed struct {
	Data   pricejson `json:"data"`
	Status string    `json:"status"`
}

func StringToByte32(str string) ([32]byte, error) {
	result := [32]byte{}
	if len(str) <= 2 {
		return result, errors.New(str + " is invalid hex")
	}
	if len(str)%2 == 1 {
		str = "0x0" + str[2:]
	}
	temp, err := hexutil.Decode(str)
	if err != nil {
		return result, err
	}
	copy(result[:32], temp)
	return result, nil
}

func (self *PriceFeed) Feed() (*Price, error) {
	r, err := StringToByte32(self.Data.R)
	if err != nil {
		return nil, err
	}
	s, err := StringToByte32(self.Data.S)
	if err != nil {
		return nil, err
	}
	return &Price{
		Block:   self.Data.Block,
		Ask:     self.Data.Ask,
		Bid:     self.Data.Bid,
		Hash:    self.Data.Hash,
		Nonce:   self.Data.Nonce,
		R:       r,
		S:       s,
		V:       self.Data.V,
		Signer:  self.Data.Signer,
		Message: self.Data.Message,
	}, nil
}

type FeedCorpus struct {
	client *http.Client
}

func (self *FeedCorpus) GetFeed() (blockNumber *big.Int, nonce *big.Int, ask1KDigix *big.Int, bid1KDigix *big.Int, v uint8, r [32]byte, s [32]byte, err error) {
	f, err := self.GetFeedFromEndpoint()
	if err != nil {
		return nil, nil, nil, nil, 0, [32]byte{}, [32]byte{}, err
	}
	return f.Block, f.Nonce, f.Ask, f.Bid, f.V, f.R, f.S, nil
}

func (self *FeedCorpus) GetFeedFromEndpoint() (*Price, error) {
	result := PriceFeed{}
	r, err := self.client.Get(ENDPOINT)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("data: %s", data)
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "success" {
		return nil, errors.New("The price feed endpoint returns unsuccessfully")
	} else {
		return result.Feed()
	}
}

func NewFeedCorpus() *FeedCorpus {
	return &FeedCorpus{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}
