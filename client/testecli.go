package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"

	"github.com/Seriyin/lab0-go/bank"
)

func main() {
	ch := make(chan int64, 200)
	n := rand.Intn(12) + 4
	ip := net.IPv4(127, 0, 0, 1)
	tcp := new(net.TCPAddr)
	tcp.IP = ip
	tcp.Port = 22556
	conn, err := net.DialTCP("tcp", nil, tcp)
	if err != nil {
		panic(err)
	}
	master, err := conn.File()
	conn.Close()
	defer master.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed opening in blocking mode %s", err)
	}
	dec := gob.NewDecoder(master)
	enc := gob.NewEncoder(master)
	for i := n; i > 0; i-- {
		conn, err := net.DialTCP("tcp", nil, tcp)
		if err != nil {
			// handle error
			panic(err)
		}
		f, err := conn.File()
		if err != nil {
			panic(err)
		}
		conn.Close()
		go spamOps(f, ch)
	}
	r := int64(0)
	for ; n > 0; n-- {
		r += <-ch
	}
	rep := bank.Reply{}
	enc.Encode(bank.Message{Op: 0, Mov: 0})
	dec.Decode(rep)
	fmt.Printf("Got %d, Expected %d\n", r, rep.Balance)
}

func spamOps(conn *os.File, ch chan int64) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	r := bank.Reply{}
	sum := int64(0)
	for i := rand.Intn(30000) + 50000; i > 0; i-- {
		enc.Encode(bank.Message{Op: 1, Mov: rand.Int63n(400) - 200})
		dec.Decode(r)
		sum += r.Balance
	}
}
