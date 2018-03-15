package main

//Account is a struct with an integer balance.
//It implements the Bank interface.
type Account struct {
	bal uint64
}

func (acc *Account) balance() uint64 {
	return acc.bal
}

func (acc *Account) movement(mov int64) bool {
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