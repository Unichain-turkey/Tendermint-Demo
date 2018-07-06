package transactoin

import (
	"crypto/sha256"
	"strconv"
)

type Transaction struct {
	Txid      int
	Hash	  [32]byte
	Name	  string
	Major	  string
	Identiy	  int
	Signature []byte//
	PubKey    []byte//of scholl

}
func (t *Transaction) MixTransaction() string{

	return strconv.Itoa(t.Identiy) +t.Name +t.Major

}

func NewTransaction(_name string, _major string,_id int) *Transaction {

	hash := sha256.Sum256([]byte(_name+_major+strconv.Itoa(_id)))

	transaction := Transaction{0, hash,_name,_major,_id,nil,nil}

	return &transaction
}
