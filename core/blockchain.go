package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type Blockchain struct {
	Logger log.Logger
	headers []*Header
	store Storage

	lock sync.RWMutex
	validator Validator
}

func NewBlockchain(l log.Logger,genesis *Block) (*Blockchain,error){
	bc:=&Blockchain{
		Logger: l,
		headers:[]*Header{},
		store:NewMemorystore(),
	}
	bc.validator=NewBlockValidator(bc)

	err:=bc.AddBlockWithoutValidation(genesis)
	return bc,err
}

func (bc *Blockchain) SetValidator(v Validator){
	bc.validator=v
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0, 1, 2 ,3] => 4 len
// [0, 1, 2 ,3] => 3 height
func (bc *Blockchain) Height() uint32 {

	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1)

}

func (bc *Blockchain) AddBlock(b *Block) error{

	if err:=bc.validator.ValidateBlock(b); err!=nil{
		return err
	}
	return bc.AddBlockWithoutValidation(b)
}



func (bc *Blockchain) AddBlockWithoutValidation(b *Block) error{

	bc.lock.Lock()
	bc.headers=append(bc.headers, b.Header)
	bc.lock.Unlock()

	// logrus.WithFields(logrus.Fields{
	// 	"height":b.Height,
	// 	"hash":b.Hash(BlockHasher{}),
	// }).Info("Adding new block")

	bc.Logger.Log("msg","adding new block","hash",b.Hash(BlockHasher{}),"height",b.Height,"transactions",len(b.Transactions))

	return bc.store.Put(b)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header,error){

	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height],nil

}