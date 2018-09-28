package transaction

/**
转账
*/
type Transfer struct {
	FromID int
	ToID   int
	Cash   int
}

/**
一次事务中可以有多个转账
*/
type Transaction struct {
	ID    int
	Trans []Transfer
}

const (
	REQUEST_START RequestType = iota
	REQUEST_PUT
	REQUEST_COMMIT
	REQUEST_UNDO
)

type RequestType int

type Request struct {
	RequestType RequestType
	Transaction *Transaction
}

func (t *Transaction) SendTransaction() {

}
