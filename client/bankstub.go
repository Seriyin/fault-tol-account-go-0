package main

import (
	"encoding/gob"
	"net"

	"../bank"
)

type bankStub struct {
	conn net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
}

func (b *bankStub) Balance() int64 {
	r := new(bank.Reply)
	b.enc.Encode(bank.Message{Op: 0, Mov: 0})
	b.dec.Decode(r)
	return r.Balance
}

func (b *bankStub) Movement(mov int64) bool {
	r := new(bank.Reply)
	b.enc.Encode(bank.Message{Op: 1, Mov: mov})
	b.dec.Decode(r)
	return r.Res
}
