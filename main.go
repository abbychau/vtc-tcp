// at as a tcp client and connect to tcp 10.20.30.60:8899
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
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
		}
	}
}