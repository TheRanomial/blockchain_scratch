package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/TheRanomial/Blockchain_golang/core"
	"github.com/TheRanomial/Blockchain_golang/crypto"
	"github.com/TheRanomial/Blockchain_golang/types"
	"github.com/go-kit/log"
)

const defaultBlockTime=5*time.Second

type ServerOpts struct{
	ID 				string
	Transport       Transport
	Logger 			log.Logger
	RPDecodeFunc 	RPDecodeFunc
	RPCProcessor 	RPCProcessor
	Transports 		[]Transport
	Blocktime 		time.Duration
	PrivateKey 		*crypto.PrivateKey
}

type Server struct {
	ServerOpts
	mempool 	*TxPool

	// peermap 	map[Netaddr]*TcpPeer
	// mu 			sync.RWMutex

	chain 		*core.Blockchain
	isValidator bool
	rpcCh 		chan RPC
	quitCh 		chan struct{}
}

func NewServer(opts ServerOpts) (*Server,error){
	if opts.Blocktime==time.Duration(0){
		opts.Blocktime=defaultBlockTime
	}

	if opts.RPDecodeFunc==nil {
		opts.RPDecodeFunc=DefaultRPCDecodeFunc
	}

	if opts.Logger==nil {
		opts.Logger=log.NewLogfmtLogger(os.Stderr)
		opts.Logger=log.With(opts.Logger,"ID",opts.ID)
	}

	chain,err:=core.NewBlockchain(opts.Logger,GenesisBlock())

	if err!=nil{
		return nil,err
	}

	s:= &Server{
		ServerOpts: opts,
		mempool: NewTxPool(1000),
		chain: chain,
		isValidator: opts.PrivateKey!=nil,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{},1),
	}

	//if we dont have any given processor option in opts
	//we use our server as default
	if s.RPCProcessor==nil{
		s.RPCProcessor=s
	}

	if s.isValidator{
		go s.ValidatorLoop()
	}

	return s,nil
}

func (s *Server) Start(){
	s.initTransport()
	free: 
		for{
			select{
			case rpc:= <-s.rpcCh:
				msg,err:=s.RPDecodeFunc(rpc)
				if err!=nil{
					s.Logger.Log("error",err)
				}
				if err:=s.RPCProcessor.ProcessMessage(msg);err!=nil{
					s.Logger.Log("error",err)
				}

			case <-s.quitCh:
				break free
			}
		}

	s.Logger.Log("msg","server is shutting down")
}

func (s *Server) ValidatorLoop(){
	ticker:=time.NewTicker(s.Blocktime)

	s.Logger.Log("msg","Starting the validator loop","Blocktime",s.Blocktime)

	for{
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *Decodedmessage) error {

	switch t:=msg.Data.(type){
		case *core.Transaction:
			return s.processTransaction(t)
		case *core.Block:
			return s.processBlock(t)
		case *GetStatusMessage:
			return s.processGetStatusMessage(msg.From, t)
		case *StatusMessage:
			return s.processStatusMessage(msg.From, t)
		// case *GetBlocksMessage:
		// 	return s.processGetBlocksMessage(msg.From, t)
		// case *BlocksMessage:
		// 	return s.processBlocksMessage(msg.From, t)
	}
	return nil
}

func (s *Server) processStatusMessage(from Netaddr,msg *StatusMessage) error{
	fmt.Printf("=>received status message from %s =>%+v\n",from,msg)

	return nil
}

func (s *Server) processGetStatusMessage(from Netaddr,data *GetStatusMessage) error{
	fmt.Printf("=>received get status message from %s =>%+v\n",from,data)

	statusMessage:=&StatusMessage{
		CurrentHeight: s.chain.Height(),
		ID: s.ID,
	}

	buf:=&bytes.Buffer{}
	if err:=gob.NewEncoder(buf).Encode(statusMessage); err!=nil{
		return err
	}

	msg:=NewMessage(MessageTypeStatus,buf.Bytes())

	return s.Transport.SendMessage(from,msg.Bytes())
}

func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		s.Logger.Log("error is there", err.Error())
		return err
	}

	go s.broadcastBlock(b)

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {

	hash:=tx.Hash(core.TxHasher{})

	if s.mempool.Contains(hash){
		return nil
	}
	
	if err := tx.Verify(); err != nil {
		return err
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	
	//s.Logger.Log("msg","adding new transaction to mempool","hash",hash,"mempoolLength",s.mempool.PendingCount())

	go s.broadcastTx(tx)
	
	s.mempool.Add(tx)

	return nil
}

func (s *Server) broadcastBlock(b *core.Block) error{
	buf:=&bytes.Buffer{}
	if err:=b.Encode(core.NewGobBlockEncoder(buf)); err!=nil{
		return err
	}

	msg:=NewMessage(MessageTypeBlock,buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcastTx(tx *core.Transaction) error {

	buf:=&bytes.Buffer{}
	if err:=tx.Encode(core.NewGobTxEncoder(buf)); err!=nil {
		return err
	}

	msg:=NewMessage(MessageTypeTx,buf.Bytes())
	
	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcast(payload []byte) error {
	for _,tr:=range s.Transports{
		if err:=tr.Broadcast(payload);err!=nil{
			return err
		}
	}
	return nil
}

func (s *Server) initTransport(){
	for _,tr:=range s.Transports{
		go func (tr Transport){
			for rpc:= range tr.Consume(){
				s.rpcCh<-rpc
			}
		}(tr)
	}
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}
	txx := s.mempool.Pending()

	block, err := core.NewBlockFromPreviousHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.mempool.ClearPending()

	go s.broadcastBlock(block)

	return nil
}

func GenesisBlock() *core.Block {
	header:=&core.Header{
		Version: 1,
		DataHash: types.Hash{},
		Height: 0,
		Timestamp: 000000,
	}

	b, _ := core.NewBlock(header, nil)

	return b
}