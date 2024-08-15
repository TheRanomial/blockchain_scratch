package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/TheRanomial/Blockchain_golang/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct {}

type TxHasher struct {}

func (BlockHasher) Hash(b *Header) types.Hash{
	header := sha256.Sum256(b.Bytes())
	return types.Hash(header)
}

//hashing all the values of a transaction
func (TxHasher) Hash(tx *Transaction) types.Hash{
	buf:=&bytes.Buffer{}

	binary.Write(buf,binary.LittleEndian,tx.Data)
	//binary.Write(buf,binary.LittleEndian,tx.To)
	//binary.Write(buf,binary.LittleEndian,tx.Value)
	//binary.Write(buf,binary.LittleEndian,tx.From)
	//binary.Write(buf,binary.LittleEndian,tx.Nonce)

	return types.Hash(sha256.Sum256(buf.Bytes()))
}