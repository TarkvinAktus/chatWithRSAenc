package main

import (
	"fmt"
	"os"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
)


func main(){
	ServerPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        fmt.Println(err.Error())
    }

	ServerPublicKey := &ServerPrivateKey.PublicKey

    CheckerPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

	CheckerPublicKey := &CheckerPrivateKey.PublicKey
	

	checkerFile, err := os.OpenFile("CheckerPrivateKey.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}else{
		PrivateX509 := x509.MarshalPKCS1PrivateKey(CheckerPrivateKey)
		if _, err := checkerFile.Write(PrivateX509); err != nil {
			checkerFile.Close() // ignore error; Write error takes precedence
			fmt.Println(err.Error())
		}
	}
	checkerFile.Close()

	checkerFile, err = os.OpenFile("CheckerPublicKey.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}else{
		PublicX509 := x509.MarshalPKCS1PublicKey(CheckerPublicKey)
		if _, err := checkerFile.Write([]byte(PublicX509)); err != nil {
			checkerFile.Close() // ignore error; Write error takes precedence
			fmt.Println(err.Error())
		}
	}
	checkerFile.Close()


	serverFile, err := os.OpenFile("ServerPrivateKey.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}else{
		PrivateX509 := x509.MarshalPKCS1PrivateKey(ServerPrivateKey)
		if _, err := serverFile.Write(PrivateX509); err != nil {
			serverFile.Close() // ignore error; Write error takes precedence
			fmt.Println(err.Error())
		}
	}
	serverFile.Close()

	serverFile, err = os.OpenFile("ServerPublicKey.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}else{
		PublicX509 := x509.MarshalPKCS1PublicKey(ServerPublicKey)
		if _, err := serverFile.Write(PublicX509); err != nil {
			serverFile.Close() // ignore error; Write error takes precedence
			fmt.Println(err.Error())
		}
	}
	serverFile.Close()











    
}