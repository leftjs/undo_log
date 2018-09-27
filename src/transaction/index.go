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

func NewTransaction() {
}
