package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/TheRanomial/Blockchain_golang/core"
	mycrypto "github.com/TheRanomial/Blockchain_golang/crypto"
	"github.com/TheRanomial/Blockchain_golang/network"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
)

func main(){
	addr1:= &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3000,
    }
	trLocal:=network.NewLocalTransport(addr1)

	addr2 := &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3002,
    }
	trRemoteA:=network.NewLocalTransport(addr2)

	addr3:= &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3003,
    }
	trRemoteB:=network.NewLocalTransport(addr3)
	
	addr4 := &net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 3004,
    }
	trRemoteC:=network.NewLocalTransport(addr4)

	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)
	trRemoteB.Connect(trRemoteA)
	trRemoteA.Connect(trLocal)

	privKey:=mycrypto.GeneratePrivateKey()
	initRemoteServers([]network.Transport{trRemoteA,trRemoteB,trRemoteC})
	go func(){
		for {
			if err:=SendTransaction(trRemoteA,trLocal.Addr());err!=nil{
				logrus.Error(err)
			}
			time.Sleep(2*time.Second)
		}
	}()

	if err:=sendGetStatusMessage(trRemoteA,trRemoteB.Addr()); err!=nil{
		log.Fatal(err)
	}

	/*go func(){
		time.Sleep(7*time.Second)

		addrLate:= &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 3005,
		}
		trLate:=network.NewLocalTransport(addrLate)
		trRemoteC.Connect(trLate)
		lateServer:=makeServer("REMOTE_LATE",nil,trLate)
		lateServer.Start()
	}()*/

	localServer:=makeServer("LOCAL",&privKey,trLocal)
	localServer.Start()
}

func sendGetStatusMessage(tr network.Transport,to network.Netaddr) error {
	var (
		getStatusMsg=new(network.GetStatusMessage)
		buf=new(bytes.Buffer)
	)

	if err:=gob.NewEncoder(buf).Encode(getStatusMsg);err!=nil{
		return err
	}

	msg:=network.NewMessage(network.MessageTypeGetStatus,buf.Bytes())

	if err:=tr.SendMessage(to,msg.Bytes());err!=nil{
		return err
	}

	return nil
}

func initRemoteServers(trs []network.Transport){
	for i:=0;i<len(trs);i++{
		id:=fmt.Sprintf("REMOTE_%d",i)
		s:=makeServer(id,nil,trs[i])
		go s.Start()
	}
}

func makeServer(id string,privKey *mycrypto.PrivateKey,tr network.Transport) *network.Server{

	opts:=network.ServerOpts{
		Transport: tr,
		PrivateKey: privKey,
		ID:id,
		Transports: []network.Transport{tr},
	}

	s,err:=network.NewServer(opts)
	if err!=nil{
		log.Fatal(err)
	}
	return s
}

func SendTransaction(tr network.Transport, to net.Addr) error{
	privKey:=mycrypto.GeneratePrivateKey()
	data:=[]byte(strconv.FormatInt(int64(rand.Intn(1000)),10))
	tx:=core.NewTransaction(data)
	tx.Sign(privKey)

	buf:=&bytes.Buffer{}

	if err:=tx.Encode(core.NewGobTxEncoder(buf)); err!=nil{
		return err
	}

	msg:=network.NewMessage(network.MessageTypeTx,buf.Bytes())

	return tr.SendMessage(to,msg.Bytes())
}





