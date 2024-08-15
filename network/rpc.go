package network

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/TheRanomial/Blockchain_golang/core"
	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3
	MessageTypeStatus    MessageType = 0x4
	MessageTypeGetStatus MessageType = 0x5
	MessageTypeBlocks    MessageType = 0x6
)

type RPC struct{
	From net.Addr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType,data []byte) *Message{
	return &Message{
		Header:t,
		Data:data,
	}
}

func (msg *Message) Bytes() []byte{
	buf:=&bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type Decodedmessage struct {
	From Netaddr
	Data any
}

type RPDecodeFunc func(RPC) (*Decodedmessage,error)

func DefaultRPCDecodeFunc(rpc RPC) (*Decodedmessage,error) {

	msg:=Message{}
	if err:=gob.NewDecoder(rpc.Payload).Decode(&msg);err!=nil{
		return nil,fmt.Errorf("failed to decode message from %s:%s",rpc.From,err)
	}

	logrus.WithFields(logrus.Fields{
		"From":rpc.From,
		"Message":msg.Header,
	}).Debug("New incoming message")

	switch msg.Header {
		case MessageTypeTx:
		tx:=new(core.Transaction)
		if err:=tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))) ; err!=nil{
			return nil,err
		}
		return &Decodedmessage{
			From:rpc.From,
			Data:tx,
		},nil

		case MessageTypeBlock:
		b:=new(core.Block)
		if err:=b.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err!=nil{
			return nil,err
		}
		return &Decodedmessage{
			From:rpc.From,
			Data:b,
		},nil

		case MessageTypeGetStatus:
		return &Decodedmessage{
			From: rpc.From,
			Data:&GetStatusMessage{},
		},nil

		case MessageTypeStatus:
		statusMessage:=new(StatusMessage)
		if err:=gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&statusMessage);err!=nil{
			return nil,err
		}
		return &Decodedmessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil

		case MessageTypeBlocks:
		getBlocksMsg:=new(GetBlocksMessage)
		if err:=gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&getBlocksMsg);err!=nil{
			return nil,err
		}
		return &Decodedmessage{
			From: rpc.From,
			Data:getBlocksMsg,
		},nil

		case MessageTypeGetBlocks:
			getBlocks := new(GetBlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocks); err != nil {
			return nil, err
		}
		return &Decodedmessage{
			From: rpc.From,
			Data: getBlocks,
		}, nil

	default:
		return nil,fmt.Errorf("invalid message header %x",msg.Header)
	}

}

type RPCProcessor interface{
	ProcessMessage(*Decodedmessage) error
}

func init() {
	gob.Register(elliptic.P256())
}