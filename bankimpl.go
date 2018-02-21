package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type reply struct {
	op      byte
	res     bool
	balance uint64
}

type message struct {
	op  byte
	mov int64
	rep chan reply
}

func main() {
	// Listen on TCP port 2000 on all available unicast and
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
	accchan := make(chan message, 100)
	go accounter(accchan)
	for {
		// Wait for a connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed connect %s", err)
		}
		f, err := conn.File()
		conn.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed opening in blocking mode %s", err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConn(f, accchan)
	}

}

func accounter(accchan chan message) {
	acc := new(Account)
	for {
		m := <-accchan
		switch m.op {
		case 0:
			m.rep <- reply{0, true, acc.balance()}
		case 1:
			m.rep <- reply{1, acc.movement(m.mov), 0}
		case 10:
			return
		default:
			m.rep <- reply{255, false, 0}
		}
	}
}

func handleConn(c *os.File, a chan message) {
	// Shut down the connection.
	defer c.Close()
	in := bufio.NewScanner(c)
	out := bufio.NewWriter(c)
	ch := make(chan reply, 50)
	for in.Scan() {
		if handleWords(ch, a, in.Text()) {
			handleReply(ch, out)
		}
	}
	if err := in.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning %s", err)
	}
}

func handleWords(ch chan reply, a chan message, text string) (res bool) {
	words := strings.Split(text, " ")
	res = true
	switch words[0] {
	case "mov":
		mov, err := strconv.ParseInt(words[1], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing mov int %s", err)
			res = false
		} else {
			a <- message{1, mov, ch}
		}
	case "bal":
		a <- message{0, 0, ch}
	default:
		res = false
	}
	return
}

func handleReply(ch chan reply, out *bufio.Writer) {
	rep := <-ch
	switch rep.op {
	case 0:
		out.WriteString(fmt.Sprintf("bal %d", rep.balance))
	case 1:
		out.WriteString(fmt.Sprintf("mov %t", rep.res))
	case 255:
		out.WriteString("unk")
	}
}
