package main

import (
"os"
"crypto/rsa"
"crypto/x509"
"crypto"
"fmt"
"net/http"
"io/ioutil"
)

func testEqBytes(a, b []byte) bool {

    // If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

func testCheckerServer(w http.ResponseWriter, r *http.Request){
	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://192.168.0.100:9009", nil)
	if err != nil {
		fmt.Println("Create request err",err.Error())
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Do request err",err.Error())
	}else{
		defer resp.Body.Close()
	
		signature, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Signature read",err.Error())
		}

		myFileKey, err := os.Open("/home/alex/Go/ChatServer/GenPrivateKeys/CheckerPublicKey.txt")
			if err != nil{
				w.Write([]byte("ERROR! NO PUBLIC KEY"))
				w.Write([]byte(err.Error()))
			}else{
				bytesData, err := ioutil.ReadAll(myFileKey)
				if err != nil {
					w.Write([]byte("ERROR! CAN'T READ FILE"))
				}

				CheckerPublicKey, err := x509.ParsePKCS1PublicKey(bytesData)
				if err != nil {
					w.Write([]byte("ERROR! CAN'T PARSE PUBLIC KEY"))
				}else{
					var opts rsa.PSSOptions  
					opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example  
					PSSmessage := []byte("Hello rsa!")  
					newhash := crypto.SHA256  
					pssh := newhash.New()  
					pssh.Write(PSSmessage)  
					hashed := pssh.Sum(nil)
					
					err = rsa.VerifyPSS(CheckerPublicKey, newhash, hashed, signature, &opts)
					if err != nil {
						w.Write([]byte("Who are U? Verify Signature failed"))  
					} else {
						w.Write([]byte("Verify Signature successful"))
					}
				}
			}
	}
}


func startServer() {
	http.HandleFunc("/", testCheckerServer )
	fmt.Println("starting server at :9009")
	http.ListenAndServe(":9009", nil)
}


func main() {
startServer()

}


