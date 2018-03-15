package bank

//Bank is a simple one account interface.
//It has a movement and a balance operation.
type Bank interface {
	Balance() int64
	Movement(int64) bool
}

//Message represents the operation to request and
//a movement to try on the account if requested.
type Message struct {
	Op  byte
	Mov int64
}

//Reply records the operation executed, current balance
//or balance moved and the result of the movement.
type Reply struct {
	Op      byte
	Res     bool
	Balance int64
}
