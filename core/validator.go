package core

import (
	"errors"
	"fmt"
)

var ErrorBlockKnown=errors.New("Block is already known")

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct{
	bc *Blockchain
}

func NewBlockValidator(b *Blockchain) *BlockValidator{
	return &BlockValidator{
		bc: b,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error{

	v.bc.Logger.Log("msg", "validating block", "height", b.Height, "hash", b.Hash(BlockHasher{}))

	if v.bc.HasBlock(b.Height){
		return ErrorBlockKnown
	}

	if b.Height!=v.bc.Height()+1 {
		return fmt.Errorf("Block with height is too high than current height")
	}

	prevheader,err:=v.bc.GetHeader(b.Height-1)

	if err!=nil{
		return fmt.Errorf("previous Header not receviable")
	}

	hash:=BlockHasher{}.Hash(prevheader)

	if hash!=b.PrevBlockHash{
		return fmt.Errorf("current and previous block hashes dont match")
	}

	if err:=b.Verify();err!=nil{
		return err
	}

	return nil
}