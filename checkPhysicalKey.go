package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
	"./constants"
)

func readDiskAndWrieSignature(writer http.ResponseWriter, request *http.Request, file *os.File) {
	
	bytesData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Read file err ", err.Error())
	}

	PrivateKey, err := x509.ParsePKCS1PrivateKey(bytesData)
	if err != nil {
		fmt.Println("Parse private key err ", err.Error())
	}

	//Create signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := []byte(constants.SignatureText)
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, PrivateKey, newhash, hashed, &opts)
	if err != nil {
		fmt.Println("Signature err ", err.Error())
	}

	fileInfo, err := os.Stat("/media/floppy/lastLogin.txt")
	if err != nil {
		fmt.Println("Login file err ", err.Error())
	} else {
		if fileInfo.Size() < 1000 {
			logFloppyFile, err := os.OpenFile("/media/floppy/lastLogin.txt", os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("Open err ", err.Error())
			} else {
				currentTime := time.Now().String()
				logFloppyFile.Write([]byte(currentTime))
			}
		} else {
			// logFloppyFile, err := os.OpenFile("/media/floppy/lastLogin.txt", os.O_RDWR, 0666)
			// if err != nil{
			// 	fmt.Println("Open err ",err.Error())
			// }else{
			// 	_ , err = logFloppyFile.Seek(0,0)
			// 	currentTime := time.Now().String()
			// 	logFloppyFile.Write([]byte(currentTime))
			// }
		}
	}
	writer.Write(signature)
}

func openDisk(writer http.ResponseWriter, request *http.Request, fileAddr string) {
	myFileKey, err := os.Open(fileAddr)
	if err != nil {
		fmt.Println("Open err ", err.Error())
	} else {
		//
		readDiskAndWrieSignature(writer, request, myFileKey)
		//
		myFileKey.Close()
	}
}

func checkPrivateKeyDisk(writer http.ResponseWriter, request *http.Request, diskAddr string, fileAddr string) {
	cmd := exec.Command("mount", diskAddr)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Mount err ", err.Error())
	} else {
		//
		openDisk(writer, request, fileAddr)
		//
		cmdUm := exec.Command("umount", diskAddr)
		errUm := cmdUm.Run()
		if errUm != nil {
			fmt.Println("Umount err ", errUm.Error())
		}
	}
}

func startServer() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		diskAddr := "/media/floppy"
		fileAddr := "/media/floppy/CheckerPrivateKey.txt"
		checkPrivateKeyDisk(writer, request, diskAddr, fileAddr)
	})

	fmt.Println("starting server at :9009")
	http.ListenAndServe(":9009", nil)
}

func main() {
	startServer()

}
