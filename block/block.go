package block

import "certificate/transaction"

var initialTime string = "2018-01-01T00:00:00.000Z"


type Block struct {
	Index           int
	Transaction		transactoin.Transaction
	Hash            string
	PrevHash        string
}

