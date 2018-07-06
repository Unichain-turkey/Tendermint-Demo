package main


import (
	"github.com/gin-gonic/gin"
	"certificate/wallet"
	"fmt"
	"certificate/transaction"
	"strconv"
)
type transactFormat struct {
	NAME 	string 	  			`json:"name""  		binding:"required"`
	MAJOR 	string 				`json:"major"  		binding:"required"`
	ID    	int					`json:"id"  		binding:"required"`
}



var mywallet * wallet.Wallet

var transactions []transactoin.Transaction


func getIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"For creating public and private keys": "/keys",
		"For signing transaction to call with post request ": "/sign ",
		"For validation transaction to call  with post request ": "/validation ",
	})
}

func getKeys(c *gin.Context) {

	mywallet=wallet.NewWallet();
	keys:=  string(mywallet.ExportRsaPrivateKey())+
			string(mywallet.ExportRsaPublicKey())

	c.String(200,"%s",keys)

}
func getAllTransactions(c *gin.Context) {

	data :=gin.H{"ItemList": "Blank",}

	for i, _tran := range transactions {
		data[strconv.Itoa(i)]=_tran
	}
	c.JSON(200, data)

}

func signTransaction(c *gin.Context) {

	var t transactFormat
	c.BindJSON(&t)

	fmt.Println(t)
	_transaction :=transactoin.NewTransaction(t.NAME,t.MAJOR,t.ID)
	signature := mywallet.SignTransaction(_transaction)
	_transaction.Signature=signature
	transactions = append(transactions,*_transaction)

	c.JSON(200, gin.H{
		"signature": signature,
	})

}
func verifyTransaction(c *gin.Context) {

	index := c.Query("index")
	fmt.Println(index)

	i, _ := strconv.Atoi(index)

	_transact := transactions[i]
	succes := mywallet.VerifyTransaction(&_transact)
	c.JSON(200, gin.H{
		"index": 	i,
		"status":succes,
		"transaction":transactions[i],

	})

}
func main() {


	router := gin.Default()
	router.GET("/", getIndex)
	router.GET("/keys", getKeys)
	router.GET("/all", getAllTransactions)
	router.POST("/sign", signTransaction)
	router.GET("/verify", verifyTransaction)


	router.Run() // listen and serve on 0.0.0.0:8080
}