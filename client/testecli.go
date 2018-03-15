package main

import (
	"fmt"
	"math/rand"
)

func main() {
	n := rand.Intn(12) + 4
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
	r := rand.Intn(30000) + 50000
	fmt.Printf("%d Spammer will spam %d iterations", index, r)
	for i := 0; i < r; i++ {
		mov := rand.Int63n(400) - 200
		if b.Movement(mov) {
			sum += mov
		} else {
			fmt.Printf("Rejected %d\n", mov)
		}
	}
	fmt.Printf("%d Spammer got %d", index, sum)
	ch <- sum
	close(ch)
}
