package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"time"

	"github.com/TheRanomial/Blockchain_golang/crypto"
	"github.com/TheRanomial/Blockchain_golang/types"
)

type Header struct {
	Version uint32
	DataHash  types.Hash
	PrevBlockHash types.Hash
	Timestamp uint64
	Height	uint32
}

type Block struct{
	*Header
	Transactions []*Transaction

	Validator  crypto.PublicKey
	Signature  *crypto.Signature

	//cache hash so that we can retreive it
	hash types.Hash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

func NewBlock(h *Header,txx []*Transaction) (*Block,error){
	return &Block{
		Header: h,
		Transactions: txx,
	},nil
}

func NewBlockFromPreviousHeader(prevHeader *Header, txx []*Transaction) (*Block,error) {

	dataHash,err:=CalculateDataHash(txx)
	if err!=nil{
		return nil,err
	}

	h:=&Header{
		Version:1,
		DataHash: dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp: uint64(time.Now().UnixNano()),
		Height:prevHeader.Height + 1,
	}

	return NewBlock(h,txx)
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
	hash, _ := CalculateDataHash(b.Transactions)
	b.DataHash = hash
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {

	sig,err:=privKey.Sign(b.Header.Bytes())
	if err!=nil{
		return err
	}
	b.Validator=privKey.PublicKey()
	b.Signature=sig

	return nil

}

func (b *Block) Verify() error{
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator,b.Header.Bytes()){
		return fmt.Errorf("block has invalid signature")
	}
	
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", BlockHasher.Hash(BlockHasher{},b.Header))
	}

	return nil
}

func (b *Block) Encode(enc Encoder[*Block]) error{
	return enc.Encode(b)
}

func (b *Block) Decode(dec Decoder[*Block]) error{
	return dec.Decode(b)
}

func CalculateDataHash(txx []*Transaction) (hash types.Hash,err error){

	buf:=&bytes.Buffer{}
	for _,tx:=range txx{
		enc := NewGobTxEncoder(buf)
		if err=enc.Encode(tx); err!=nil{
			return 
		}
	}

	hash = sha256.Sum256(buf.Bytes())
	return 
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}






