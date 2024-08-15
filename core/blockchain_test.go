package core

import (
	"fmt"
	"testing"

	"github.com/TheRanomial/Blockchain_golang/types"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_Blockchain(t *testing.T){
	bc:=NewBlockchainWithGenesis(t)

	assert.NotNil(t,bc.validator)
	fmt.Println(bc.Height())
	assert.Equal(t,bc.Height(),uint32(0))
}

func Test_AddBlock(t *testing.T){
	bc:=NewBlockchainWithGenesis(t)

	lenBlocks:=10
	for i:=0;i<10;i++{
		block:=randomBlock(t, uint32(i+1), getPrevBlockHash(t,bc,uint32(i+1)))
		assert.Nil(t,bc.AddBlock(block))
	}

	assert.Equal(t,bc.Height(),uint32(lenBlocks))
	assert.Equal(t,len(bc.headers),lenBlocks+1)
}

func TestAddBlockToHigh(t *testing.T) {
	bc := NewBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}

func getPrevBlockHash(t *testing.T, bc *Blockchain,height uint32) types.Hash {

	prevHeader,err:=bc.GetHeader(height-1)
	assert.Nil(t,err)
	return BlockHasher{}.Hash(prevHeader)
}

func NewBlockchainWithGenesis(t *testing.T) *Blockchain{
	bc,err:=NewBlockchain(log.NewNopLogger(),randomBlock(t, 0, types.Hash{}))
	assert.Nil(t,err)
	return bc
}
