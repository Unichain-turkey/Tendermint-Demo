package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	"strconv"
	"strings"
	"time"
	"certificate/transaction"
	"encoding/json"
    "net/http"
	"io/ioutil"
)

var initialTime string = "2018-01-01T00:00:00.000Z"

type Data struct {
	Id   int
	Key  int
	Date string
}
type Block struct {
	Index           int
	Timestamp       string
	Transaction		transaction.Transaction
	Hash            string
	PrevHash        string

}

type Certificate struct {
	types.BaseApplication
	Blockchain []Block
	Height     int64
	AppHash    []byte
}




var pendingTransaction [] transaction.Transaction;

func NewCertificateApplication() *Certificate {
	var _blockchain []Block

	t, _ := time.Parse(time.RFC3339, initialTime)

	h := sha256.New()
	h.Write([]byte(string(0) + t.String()))


	genesisBlock := Block{0, t.String(), transaction.Transaction{}, hex.EncodeToString(h.Sum(nil)),""}

	_blockchain = append(_blockchain, genesisBlock)

	return &Certificate{Height: 0, Blockchain: _blockchain, AppHash: nil}
}
func (app *Certificate) Info(req types.RequestInfo) types.ResponseInfo {
	fmt.Println("Info: ", app.Height, " size ", len(app.Blockchain))
	return types.ResponseInfo{Data: cmn.Fmt("{\"hashes\":%v,\"height\":%v}", app.AppHash, app.Height)}
}
func (app *Certificate) DeliverTx(tx []byte) types.ResponseDeliverTx {

	if len(tx) != 1 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Incomplete Argument List")}
	}
	index, _ := strconv.Atoi(string(tx))

	tras :=app.getTransaction(index)
	fmt.Println(tras)


	newBlock, err := generateBlock(app.Blockchain[len(app.Blockchain)-1],tras)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeBadNonce,
			Log:  fmt.Sprintf("Invalid block type.")}
	}
	if isBlockValid(newBlock, app.Blockchain[len(app.Blockchain)-1]) {
		newBlockchain := append(app.Blockchain, newBlock)

		if len(newBlockchain) > len(app.Blockchain) {
			app.Blockchain = newBlockchain
		}
		spew.Dump(app.Blockchain)
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *Certificate) CheckTx(tx []byte) types.ResponseCheckTx {
	fmt.Println("in CheckTx ", tx)

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func byteToTransaction(tx []byte) transaction.Transaction{
	var _trans transaction.Transaction
	json.Unmarshal(tx, &_trans)
	return _trans
}

func (app *Certificate) Commit() (resp types.ResponseCommit) {

	if app.Height == 0 {
		return types.ResponseCommit{}
	}
	app.Height += 1
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, int64(app.Height))

	app.AppHash = appHash

	fmt.Println("in Commit ", appHash)

	return types.ResponseCommit{Data: appHash}
}

func (app *Certificate) Query(reqQuery types.RequestQuery) types.ResponseQuery {

	query := strings.Split(string(reqQuery.Data), "=")
	fmt.Println("in Query Data ", query)

	switch query[0] {
	case "hash":
		for _, v := range app.Blockchain {
			if string(v.Transaction.Hash[:]) == query[1] {
				return types.ResponseQuery{Log: "exists"}
			}
		}
		return types.ResponseQuery{Log: cmn.Fmt("Not found This hash", query[1])}

	default:
		return types.ResponseQuery{Log: cmn.Fmt("Invalid query path", reqQuery.Path)}
	}
}

func calculateHash(block Block) string {

	record := string(block.Index) + block.Timestamp + string(block.Transaction.Hash[:]) + string(block.Hash) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)

}
func generateBlock(oldBlock Block, transact transaction.Transaction) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transaction = transact

	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}
func (app *Certificate) getTransaction(index int) transaction.Transaction {
	resp, err := http.Get("http://localhost:8080/getOne?index="+string(index))
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)



	arr := []byte(string(body))

	var _trans transaction.Transaction
	json.Unmarshal(arr, &_trans)
	return _trans

}

func main() {

	fmt.Println("Started My app")
	app := NewCertificateApplication()

	srv := server.NewSocketServer("tcp://0.0.0.0:46658", app)

	if err := srv.Start(); err != nil {
		fmt.Println(err)
	}
	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}
