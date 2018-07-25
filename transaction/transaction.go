package transaction

import (
	"crypto/sha256"
	"strconv"
	"fmt"
	"encoding/json"
	"time"
)

type Transaction struct {
	Txid      int
	Hash	  [32]byte
	Name	  string
	Major	  string
	Identiy	  int
	Signature []byte//
	PubKey    []byte//of scholl
	Timestamp string

}
func (t *Transaction) MixTransaction() string{

	return strconv.Itoa(t.Identiy) +t.Name +t.Major

}

func NewTransaction(_name string, _major string,_id int) *Transaction {

	hash := sha256.Sum256([]byte(_name+_major+strconv.Itoa(_id)))
	_timestamp:= time.Now().Format(time.RFC850)

	transaction := Transaction{0, hash,_name,_major,_id,nil,nil,_timestamp}

	return &transaction
}

func test()  {
	// Create a struct and write it.
	t := *NewTransaction("Fatih","Computer Science",150113082)
	//x := Transaction{1, [32]byte{43, 1, 0},"a","b",1,[4]byte{43, 1, 0},[4]byte{43, 1, 0},"a"}
	b, _ := json.Marshal(t)
	// Convert bytes to string.
	fmt.Println(b)
	s := string(b)
	fmt.Println(s)
	var languages Transaction
	json.Unmarshal(b, &languages)
	fmt.Println(languages)

}
