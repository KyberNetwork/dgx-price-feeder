package dgxpricing

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Blockchain interface {
	// status should be:
	// 1. "": pending
	// 2. lost: not found from the node
	// 3. mined: successfully
	// 4. failed: threw/reverted
	TxStatus(common.Hash) (status string, blockno uint64, err error)
}

// StatusMonitor is not thread safe
type StatusMonitor struct {
	// txs always has at least 1 tx
	txs []*types.Transaction
}

func (self *StatusMonitor) GetOneStatus(tx *types.Transaction, bc Blockchain, data *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	// we ignore the error here because we will consider the status as "" in case there is error
	status, _, _ := bc.TxStatus(tx.Hash())
	if status == "" {
		// convert "" to "pending"
		status = "pending"
	}
	data.Store(tx.Hash().Hex(), status)
}

func (self *StatusMonitor) ConcurrentlyGetStatus(bc Blockchain) map[string]string {
	data := sync.Map{}
	wg := sync.WaitGroup{}
	for _, tx := range self.txs {
		wg.Add(1)
		go self.GetOneStatus(tx, bc, &data, &wg)
	}
	wg.Wait()
	result := map[string]string{}
	data.Range(func(k, v interface{}) bool {
		result[k.(string)] = v.(string)
		return true
	})
	return result
}

func (self *StatusMonitor) GetTxByHash(hash string) *types.Transaction {
	for _, tx := range self.txs {
		if tx.Hash().Hex() == hash {
			return tx
		}
	}
	return nil
}

// GetStatus query statuses of all txs, it returns:
// 1. mined: if one of the txs is mined
// 2. failed: if one of the txs is failed
// 3. lost: if not in the case of 1 nor 2 and the last tx is not found
// 4. pending: if not in the case of 1 nor 2 nor 3 and the last tx is pending
func (self *StatusMonitor) GetStatus(bc Blockchain) (st string, tx *types.Transaction, err error) {
	statuses := self.ConcurrentlyGetStatus(bc)
	// check if any txs is mined
	for hash, status := range statuses {
		if status == "mined" {
			return status, self.GetTxByHash(hash), nil
		}
	}
	// check if any txs is failed
	for hash, status := range statuses {
		if status == "failed" {
			return status, self.GetTxByHash(hash), nil
		}
	}
	// check the last one to see if it's lost
	lastHash := self.txs[len(self.txs)-1].Hash().Hex()
	lastStatus, found := statuses[lastHash]
	if !found {
		return "", nil, errors.New("Couldn't get status of the txs")
	}
	return lastStatus, self.txs[len(self.txs)-1], nil
}

func (self *StatusMonitor) PushTx(tx *types.Transaction) error {
	if tx.GasPrice().Cmp(self.txs[len(self.txs)-1].GasPrice()) != 1 {
		return errors.New("you must push tx with higher gas price than the last one")
	}
	self.txs = append(self.txs, tx)
	return nil
}

func NewStatusMonitor(initTx *types.Transaction) *StatusMonitor {
	if initTx == nil {
		panic("initTx must not be nil")
	}
	return &StatusMonitor{
		[]*types.Transaction{initTx},
	}
}
