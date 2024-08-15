package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Keypair_Success(t *testing.T){

	privkey:=GeneratePrivateKey()
	publicKey:=privkey.PublicKey()

	data:=[]byte("Hello world")
	sig,err:=privkey.Sign(data)
	assert.Nil(t, err)

	assert.True(t,sig.Verify(publicKey,data))
}

func Test_Keypair_failure(t *testing.T){

	privKey := GeneratePrivateKey()
	publicKey := privKey.PublicKey()
	msg := []byte("hello world")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	otherPrivKey := GeneratePrivateKey()
	otherPublicKey := otherPrivKey.PublicKey()

	sig2,_:=otherPrivKey.Sign([]byte("xxxxxx"))

	assert.True(t, sig.Verify(publicKey, msg))
	assert.True(t, sig2.Verify(otherPublicKey, []byte("xxxxxx")))

}