// at as a tcp client and connect to tcp 10.20.30.60:8899
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

var currentNfcId string
var lastTimestamp int64

func strContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

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

			if err != nil {
				fmt.Println("http request failed")
				continue
			}

			fmt.Println(string(body))

			// check if string(body) contains "does not"
			if strContains(string(body), "does not") {
				fmt.Println("nfcId not found")
				currentNfcId = ""
				continue
			} else {
				lastTimestamp = time.Now().Unix()
				fmt.Println("nfcId found")

				currentNfcId = string(buf[:n])
			}

			resp.Body.Close()
		}
	}
}

func main() {

	//start nfc loop
	go nfcLoop()

	//clear currentNfcId every 10 seconds
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if time.Now().Unix()-lastTimestamp > 30 {
				currentNfcId = ""
			}
		}
	}()

	// http server
	http.HandleFunc("/nfcId", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, currentNfcId)
	})

	http.ListenAndServe(":50011", nil)

}
