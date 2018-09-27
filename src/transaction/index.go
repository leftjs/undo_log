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
	TransactionID int
	trans         []Transfer
}

const (
	REQUEST_START = iota
	REQUEST_PUT
	REQUEST_COMMIT
)

type RequestType int

type Request struct {
	requestType RequestType
	transfer    Transfer
}

func NewTransaction(trans []Transfer) *Transaction {
	return &Transaction{
		trans: trans,
	}
}
