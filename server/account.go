package main

//Account is a struct with an integer balance.
//It implements the Bank interface.
type account struct {
	bal uint64
}

func (acc *account) Balance() uint64 {
	return acc.bal
}

func (acc *account) Movement(mov int64) bool {
	res := true
	if mov > 0 {
		acc.bal += uint64(mov)
	} else {
		if uint64(-mov) > acc.bal {
			res = false
		} else {
			acc.bal -= uint64(-mov)
		}
	}
	return res
}
