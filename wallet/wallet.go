package wallet

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"certificate/transaction"
)


type Wallet struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}


func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func newKeyPair() (*rsa.PrivateKey,*rsa.PublicKey){
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
		}
	return privateKey, &privateKey.PublicKey
}


func (w * Wallet) SignTransaction (t * transactoin.Transaction) {

	messageBytes := bytes.NewBufferString(t.MixTransaction())
	hash := sha512.New()
	hash.Write(messageBytes.Bytes())
	digest := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, w.PrivateKey, crypto.SHA512, digest)

	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
	} else{
		t.Signature=signature
		t.PubKey=w.ExportRsaPublicKey()
	}

}

func  (w * Wallet) VerifyTransaction(t * transactoin.Transaction) bool {

	messageBytes := bytes.NewBufferString(t.MixTransaction())
	hash := sha512.New()
	hash.Write(messageBytes.Bytes())
	digest := hash.Sum(nil)


	err := rsa.VerifyPKCS1v15(w.PublicKey, crypto.SHA512, digest, t.Signature)
	if err != nil {
		fmt.Printf("rsa.VerifyPKCS1v15 error: %V\n", err)
		return false
	}

	fmt.Println("Signature good!")
	return true
}
func(w * Wallet) ExportRsaPrivateKey() []byte {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(w.PrivateKey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return privkey_pem
}


func  (w * Wallet) ExportRsaPublicKey() ([]byte) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(w.PublicKey)
	if err != nil {
		return nil
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return pubkey_pem
}
func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}


func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}


func test() {

	fmt.Println("Started My app")



	myWallet :=NewWallet()
	transact :=transactoin.NewTransaction("Fatih","Computer Science",150113082)

	myWallet.SignTransaction(transact)
	fmt.Println(myWallet.VerifyTransaction(transact))



}
