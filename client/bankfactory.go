package main

import (
	"encoding/gob"
	"net"

	"../bank"
)

type bankFactory struct {
	tcp   *net.TCPAddr
	conns []net.Conn
}

//NewBankFactory returns a bankFactory with default network capabilities.
//Address defaults to "127.0.0.1:22556" of a BankServer.
func NewBankFactory() (bf *bankFactory) {
	ip := net.IPv4(127, 0, 0, 1)
	tcp := new(net.TCPAddr)
	tcp.IP = ip
	tcp.Port = 22556
	conns := make([]net.Conn, 0, 16)
	return &bankFactory{tcp, conns}
}

//NewBank builds a new Bank instance and returns it.
func (bf *bankFactory) NewBank() bank.Bank {
	master, err := net.DialTCP("tcp", nil, bf.tcp)
	if err != nil {
		panic(err)
	}
	enc := gob.NewEncoder(master)
	dec := gob.NewDecoder(master)
	bf.conns = append(bf.conns, master)
	return &bankStub{conn: master, enc: enc, dec: dec}
}

// CloseBanks will close all banks returned via NewBank.
// Must be called when all Banks are no longer necessary.
func (bf *bankFactory) CloseBanks() {
	for _, c := range bf.conns {
		c.Close()
	}
}
