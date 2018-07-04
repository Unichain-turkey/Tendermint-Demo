package main

import (
	"github.com/gin-gonic/gin"
	"certificate/crypto"
	"fmt"
	"strings"
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


func getIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"For taking public and private keys": "/keys",
		"For signing transaction to call with post request ": "/sign ",
		"For validation transaction to call  with post request ": "/validation ",

	})



}
func getKeys(c *gin.Context) {

	privateKey:=crypto.GetPrivateKey()
	publicKey:=crypto.GetPublicKeyFromPrivateKey(privateKey)


	/*
	file, err := os.Create("keys.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	*/


	keys:=  crypto.ExportRsaPrivateKeyAsPemStr(privateKey)+
			crypto.ExportRsaPublicKeyAsPemStr(publicKey)

	//file.WriteString(keys)

	c.String(200,"%s",keys)


}
func signData(c *gin.Context) {

	var s signFormat
	c.BindJSON(&s)
	//data := s.DATA
	strprivate := strings.Replace(s.PRIVATEKEY," ", "",-1)
	fmt.Println(strprivate)
	_,er := crypto.ParseRsaPrivateKeyFromPemStr(strprivate)
	fmt.Println(er)
	c.JSON(200, gin.H{
		"signature": "as",

	})

}
func verifyData(c *gin.Context) {


	fmt.Println(c.Request.GetBody)
	fmt.Println(c.Request.Form)
	fmt.Println(c.Request.Header)
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