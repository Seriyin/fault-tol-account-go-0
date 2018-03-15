package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Seriyin/lab0-go/bank"
)

type mov struct {
	mov int64
	rep chan bank.Reply
}

type bal struct {
	rep chan bank.Reply
}

func main() {
	// Listen on TCP port 22556 on all available unicast and
	// anycast IP addresses of the local system.
	ip := net.IPv4(127, 0, 0, 1)
	tcp := new(net.TCPAddr)
	tcp.IP = ip
	tcp.Port = 22556
	l, err := net.ListenTCP("tcp", tcp)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	accchan := make(chan interface{}, 100)
	go accounter(accchan)
	for {
		// Wait for a connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed connect %s", err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConn(conn, accchan)
	}

}

func accounter(accchan chan interface{}) {
	var acc bank.Bank
	acc = new(account)
	for {
		m := <-accchan
		switch m := m.(type) {
		case bal:
			m.rep <- bank.Reply{Op: 0, Res: true, Balance: acc.Balance()}
		case mov:
			m.rep <- bank.Reply{Op: 1, Res: acc.Movement(m.mov), Balance: 0}
		default:
			fmt.Fprintf(os.Stderr, "Unrecognized message")
		}
	}
}

func handleConn(c net.Conn, a chan interface{}) {
	// Shut down the connection.
	defer c.Close()
	fmt.Println("Handling Connection")
	dec := gob.NewDecoder(c)
	enc := gob.NewEncoder(c)
	ch := make(chan bank.Reply, 50)
	mes := new(bank.Message)
	for err := dec.Decode(mes); err != io.EOF; {
		if err == nil && handleMessage(ch, a, mes) {
			sendReply(ch, enc)
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding: %s", err)
		}
	}
}

func handleMessage(ch chan bank.Reply, a chan interface{}, mes *bank.Message) (res bool) {
	res = true
	switch mes.Op {
	case 1:
		a <- mov{mes.Mov, ch}
	case 0:
		a <- bal{ch}
	default:
		res = false
	}
	return
}

func sendReply(ch chan bank.Reply, out *gob.Encoder) {
	rep := <-ch
	out.Encode(rep)
}
