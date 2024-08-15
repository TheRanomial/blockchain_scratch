package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

type Encoder[T any] interface{
	Encode(T) error
}

type Decoder[T any] interface{
	Decode(T) error
}

type GobTxEncoder struct{
	w io.Writer
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w,
	}
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{
		r:r,
	}
}

func (g *GobTxEncoder) Encode(tx *Transaction) error{
	return gob.NewEncoder(g.w).Encode(tx)
}

func (g *GobTxDecoder) Decode(tx *Transaction) error{
	return gob.NewDecoder(g.r).Decode(tx)
}

type GobBlockEncoder struct {
	w io.Writer
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

func (dec *GobBlockDecoder) Decode(b *Block) error{
	return gob.NewDecoder(dec.r).Decode(b)
}

func init(){
	gob.Register(elliptic.P256())
}