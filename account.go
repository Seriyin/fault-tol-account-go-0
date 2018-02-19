package account

type Interface interface{
	balance() uint64
	movement(int64) bool
}

type Account struct {
	balance int64
}

func (acc Account *) balance() uint64 {
	return acc.balance
}

func (acc Account *) movement(mov int64) bool {
	res := true
	if mov > 0 {
		acc.balance += mov
	}
	else {
		if -mov > acc.balance {
			res = false
		}
		else {
			acc.balance -= -mov
		}
	}
	return res
}
