package main

import (
	"net"
	"bufio"
	"io"
	"fmt"
	"os"
	"account"
	"strings"
	"strconv"
)

type reply struct {
	op byte
	res bool
	balance int64
}

type message struct {
	op byte
	mov int64
	rep chan reply
}


func main() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", "127.0.0.1:22556")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	accchan := make(chan message, 100)
	go func accounter(accchan)
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed connect %s", err.String())
		}
		f, err := conn.File()
		conn.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed opening in blocking mode %s", err.String())
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
		m := <- accchan
		switch m.op {
		case 0: 
			m.rep <- reply{0,true,acc.balance()}
		case 1:
			m.rep <- reply{1,acc.movement(m.mov),0}
		case 10:
			return
		default:
			m.rep <- reply{255,false,0}
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
		if handleWords(a, in.Text()) {
			handleReply(ch, out)
		}
	}
	if err := f.Err(); err != nil {
		fmt.Fprintf("Error scanning %s", err.String())
	}
}

func handleWords(a chan message, text string Text) (res bool) {
	words := string.split(text," ")
	res = true
	switch words[0] {
	case "mov" :
		mov, err := strconv.ParseInt(words[1])
		if err != nil {
			fmt.Fprintf("Error parsing mov int %s",err.String())
			res = false
		}
		else {
			a <- message{1,mov,ch}
		}
	case "bal" :
		a <- message{0,0,ch}
	default:
		res = false
	}
	return
}

func handleReply(ch chan reply, out *Writer) {
	rep := <- ch
	switch rep.op {
	case 0:
		out.WriteString(fmt.Sprintln("bal%d",rep.balance))
	case 1:
		out.WriteString(fmt.Sprintln("mov%t",rep.res))
	case 255:
		out.WriteString("unk")
	}
}
