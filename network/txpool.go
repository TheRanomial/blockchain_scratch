package network

import (
	"sync"

	"github.com/TheRanomial/Blockchain_golang/core"
	"github.com/TheRanomial/Blockchain_golang/types"
)

type TxPool struct {
	all     *TxSortedMap
	pending *TxSortedMap

	// The maxLength of the total pool of transactions.
	// When the pool is full we will prune the oldest transaction.
	maxLength int
}

func NewTxPool(maxLength int) *TxPool{
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

func (p *TxPool) Add(tx *core.Transaction) {
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	if !p.all.Contains(tx.Hash(core.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

func (t *TxPool) Contains(hash types.Hash) bool{
	return t.all.Contains(hash)
}

func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.txx.Data
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}


type TxSortedMap struct {
	lock sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txx *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap{
	return &TxSortedMap{
		lookup :make(map[types.Hash]*core.Transaction),
		txx: types.NewList[*core.Transaction](),
	}
}

func (t *TxSortedMap) First() *core.Transaction{
	t.lock.RLock()
	defer t.lock.RUnlock()

	tx:=t.txx.Get(0)
	return t.lookup[tx.Hash(core.TxHasher{})]
}

func (t *TxSortedMap) Get(h types.Hash) *core.Transaction{
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

func (t *TxSortedMap) Add(tx *core.Transaction){

	hash:=tx.Hash(core.TxHasher{})

	t.lock.RLock()
	defer t.lock.RUnlock()

	_,ok:=t.lookup[hash]
	if !ok{
		t.lookup[hash]=tx
		t.txx.Insert(tx)
	}
}

func (t *TxSortedMap) Remove(h types.Hash){
	t.lock.RLock()
	defer t.lock.RUnlock()

	t.txx.Remove(t.lookup[h])
	delete(t.lookup,h)
}

func (t *TxSortedMap) Count() int{
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

func (t *TxSortedMap) Contains(h types.Hash) bool{
	t.lock.RLock()
	defer t.lock.RUnlock()

    _, ok := t.lookup[h]
	return ok
}

func (t *TxSortedMap) Clear(){
	t.lock.RLock()
	defer t.lock.RUnlock()

	t.lookup=make(map[types.Hash]*core.Transaction)
	t.txx.Clear()
}

