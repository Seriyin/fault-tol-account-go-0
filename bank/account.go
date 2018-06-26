package bank

//Account is a struct with an integer balance.
//It implements the Bank interface.
type account struct {
	bal int64
}

func NewAccount() *account {
	return &account{}
}

func (acc *account) Balance() int64 {
	return acc.bal
}

func (acc *account) Movement(mov int64) bool {
	res := true
	if mov > 0 {
		acc.bal += mov
	} else if -mov <= acc.bal {
		acc.bal += mov
	} else {
		res = false
	}
	return res
}
