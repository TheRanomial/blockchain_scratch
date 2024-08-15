package network

import "net"

type TcpPeer struct {
	conn net.Conn
	outgoing bool
}

func (t *TcpPeer) Send(b []byte) error {
	_,err:=t.conn.Write(b)

	return err
}

func(p *TcpPeer) readLoop(rpcch chan RPC){
	
}