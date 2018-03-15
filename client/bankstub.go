package main

import (
	"encoding/gob"
	"net"

	"github.com/Seriyin/lab0-go/bank"
)

type bankStub struct {
	conn net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
	r    *bank.Reply
}

func (b *bankStub) Balance() uint64 {
	b.enc.Encode(bank.Message{Op: 0, Mov: 0})
	b.dec.Decode(b.r)
	return uint64(b.r.Balance)
}

func (b *bankStub) Movement(mov int64) bool {
	b.enc.Encode(bank.Message{Op: 1, Mov: mov})
	b.dec.Decode(b.r)
	return b.r.Res
}
