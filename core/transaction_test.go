package core

import (
	"bytes"
	"testing"

	"github.com/TheRanomial/Blockchain_golang/crypto"
	"github.com/TheRanomial/Blockchain_golang/types"
	"github.com/stretchr/testify/assert"
)


func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func Test_TxEncoding_Decoding(t *testing.T){
	tx:=randomTxWithSignature(t)
	buf:=&bytes.Buffer{}

	assert.Nil(t,tx.Encode(NewGobTxEncoder(buf)))
	tx.hash = types.Hash{}

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}