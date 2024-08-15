package core

import (
	"fmt"
	//"math/rand"

	"github.com/TheRanomial/Blockchain_golang/crypto"
	"github.com/TheRanomial/Blockchain_golang/types"
)

type Transaction struct{

	Data []byte
	//To  crypto.PublicKey
	//Value uint64
	From crypto.PublicKey
	Signature  *crypto.Signature
	Nonce int64

	hash types.Hash

	FirstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
		//Nonce: rand.Int63n(1000000000000000),
	}
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error{

	hasher := TxHasher{}
	tx.hash=hasher.Hash(tx)
	
	sig,err:=privKey.Sign(tx.hash.ToSlice())
	if err!=nil {
		return err
	}

	tx.From=privKey.PublicKey()
	tx.Signature=sig

	return nil
}

func (tx *Transaction) Verify() error {

	if tx.Signature==nil{
		return fmt.Errorf("transaction has no signature")
	}

	hash := tx.Hash(TxHasher{})
	if !tx.Signature.Verify(tx.From, hash.ToSlice()) {
		return fmt.Errorf("invalid transaction signature")
	}
	return nil
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash{
	if tx.hash.IsZero(){
		tx.hash=hasher.Hash(tx)
	}
	return tx.hash
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

func (tx *Transaction) SetFirstSeen(t int64){
	tx.FirstSeen=t
}

func (tx *Transaction) GetFirstSeen() int64 {
	return tx.FirstSeen
}


