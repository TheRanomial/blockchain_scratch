package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"

	"github.com/TheRanomial/Blockchain_golang/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type Signature struct {
	R *big.Int
	S *big.Int
}

//sigining the signature
func (k PrivateKey) Sign(data []byte) (*Signature,error) {

	r,s,err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil,err
	}

	return &Signature{
		R:r,
		S:s,
	},nil
}

//verifying the Signature
func (sig Signature) Verify(pubKey PublicKey,data []byte) bool{

	x,y:=elliptic.UnmarshalCompressed(elliptic.P256(),pubKey)

	key:=&ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X: x,
		Y:y,
	}
	return ecdsa.Verify(key,data,sig.R,sig.S)
}

func NewPrivateKeyFromReader(r io.Reader) PrivateKey{

	privateK, err := ecdsa.GenerateKey(elliptic.P256(),r)
	if err!=nil{
		panic(err)
	}
	return PrivateKey{
		key:privateK,
	}
}

func GeneratePrivateKey() PrivateKey {
	return NewPrivateKeyFromReader(rand.Reader)
}

type PublicKey []byte

//saving space by shortening
func (p PrivateKey) PublicKey() PublicKey{
	return elliptic.MarshalCompressed(p.key.PublicKey,p.key.PublicKey.X,p.key.PublicKey.Y)
}

func (p PublicKey) String() string {
	return hex.EncodeToString(p)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k)

	return types.AddressFromBytes(h[len(h)-20:])
}

func (sig Signature) String() string {
	b:=append(sig.R.Bytes(),sig.S.Bytes()...)

	return hex.EncodeToString(b)
}










