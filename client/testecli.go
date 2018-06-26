package main

import (
	"fmt"
	cr "crypto/rand"
	"math/rand"
	"encoding/binary"
)

var rnd *rand.Rand

func main() {
	var seed int64
	binary.Read(cr.Reader, binary.LittleEndian, &seed)
	rnd = rand.New(rand.NewSource(seed))
	n := rnd.Intn(8)+4
	ch := make([]chan int64, n, n)
	for i := range ch {
		ch[i] = make(chan int64, 1)
	}
	bf := NewBankFactory()
	b := bf.NewBank()
	for i := 0; i < n; i++ {
		go spamOps(bf, ch[i], i)
	}
	r := int64(0)
	for _, c := range ch {
		for i := range c {
			r += i
		}
	}
	fmt.Printf("Got %d, Expected %d\n", r, b.Balance())
	bf.CloseBanks()
}

func spamOps(bf *bankFactory, ch chan int64, index int) {
	b := bf.NewBank()
	sum := int64(0)
	r := rnd.Intn(30000) + 50000
	fmt.Printf("%d Spammer will spam %d iterations\n", index, r)
	for i := 0; i < r; i++ {
		mov := rnd.Int63n(400) - 200
		if b.Movement(mov) {
			sum += mov
		} else {
			fmt.Printf("Rejected %d Movement of: %d\n", i, mov)
		}
	}
	fmt.Printf("%d Spammer got %d\n", index, sum)
	ch <- sum
	close(ch)
}
