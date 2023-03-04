// at as a tcp client and connect to tcp 10.20.30.60:8899
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var currentNfcId string

func nfcLoop() {
	for {
		conn, err := net.Dial("tcp", "10.20.30.60:8899")
		if err != nil {
			//sleep 5 seconds
			fmt.Println("Connection failed, sleep 5 seconds")
			time.Sleep(5 * time.Second)
			conn.Close()
			continue

		}

		fmt.Println("Connection established")

		// read from connection
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Read failed, sleep 5 seconds")
				time.Sleep(5 * time.Second)
				conn.Close()
				break
			}
			fmt.Println(string(buf[:n]))
			//do http request to http://localhost:3000
			resp, err2 := http.Get("http://localhost:53000/vtc-api/nfcAuthentication?nfcId=" + string(buf[:n]))
			if err2 != nil {
				fmt.Println("http request failed")
				continue
			}

			body, err := io.ReadAll(resp.Body)
			fmt.Println(string(body))
			currentNfcId = string(body)
			resp.Body.Close()
		}
	}
}

func main() {

	//start nfc loop
	go nfcLoop()

	// http server
	http.HandleFunc("/nfcId", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, currentNfcId)
	})

	http.ListenAndServe(":50011", nil)

}
