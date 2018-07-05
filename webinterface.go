package main


import (
	"github.com/gin-gonic/gin"
	"certificate/wallet"
	"fmt"
	"os"
	"log"

)
type signFormat struct {
	DATA string 	  		`json:"data""  		binding:"required"`
	PRIVATEKEY string 		`json:"privateKey"  binding:"required"`
}
type veriyfFormat struct {
	DATA string   			`json:"data" 	  	binding:"required"`
	SIGNATURE string    	`json:"signature" 	binding:"required"`
	PUBLICKEY string 		`json:"publicKey" 	binding:"required"`
}
var mywallet wallet.Wallet

func getIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"For creating public and private keys": "/keys",
		"For signing transaction to call with post request ": "/sign ",
		"For validation transaction to call  with post request ": "/validation ",
	})
}

func getKeys(c *gin.Context) {

	mywallet=wallet.NewWallet();


	keys:=  crypto.ExportRsaPrivateKeyAsPemStr(privateKey)+
			crypto.ExportRsaPublicKeyAsPemStr(publicKey)


	c.String(200,"%s",keys)


}
func signData(c *gin.Context) {

	var s signFormat
	c.BindJSON(&s)
	data := s.DATA
	strprivate :=s.PRIVATEKEY
	privateKey,err := crypto.ParseRsaPrivateKeyFromPemStr(strprivate)

	//signature := crypto.SignTransaction(privateKey,data)
	fmt.Println(privateKey)
	fmt.Println(err)
	fmt.Println(data)
	c.JSON(200, gin.H{
		"signature": "sd",

	})

}
func verifyData(c *gin.Context) {



	data := c.PostForm("data")
	public := c.PostForm("publicKey")
	signature := c.PostForm("signature")

	c.JSON(200, gin.H{
		"private": data,
		"public ": public,
		"signature ": signature,

	})

}
func main() {

	router := gin.Default()
	router.GET("/", getIndex)
	router.GET("/keys", getKeys)
	router.POST("/signData", signData)
	router.POST("/verifyData", verifyData)


	router.Run() // listen and serve on 0.0.0.0:8080
}