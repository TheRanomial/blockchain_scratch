package network

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T){

	addr1 := &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3000,
    }
	trLocal:=NewLocalTransport(addr1)

	addr2 := &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3001,
    }
	trRemote:=NewLocalTransport(addr2)

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	assert.Equal(t,trLocal.peers[trRemote.Addr()],trRemote)
	assert.Equal(t,trRemote.peers[trLocal.Addr()],trLocal)

}

func Test_Broadcast(t *testing.T){

	addr1 := &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3000,
    }
	tra:=NewLocalTransport(addr1)

	addr2 := &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3001,
    }
	trb:=NewLocalTransport(addr2)

	addr3:= &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3002,
    }
	trc:=NewLocalTransport(addr3)

	tra.Connect(trb)
	tra.Connect(trc) 

	msg:=[]byte("hello")
	assert.Nil(t,tra.Broadcast(msg))

	rpcb:=<-trb.Consume()
	b,err:=io.ReadAll(rpcb.Payload)
	assert.Nil(t,err)
	assert.Equal(t,b,msg)

	rpcc:=<-trc.Consume()
	c,err:=io.ReadAll(rpcc.Payload)
	assert.Nil(t,err)
	assert.Equal(t,c,msg)
}

