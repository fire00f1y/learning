package main

import (
	"net"
	"fmt"
	"strconv"
	"time"
	"os"
)

func main() {
	StartUdpServer()
	//RunUdpClient()
}

func RunUdpClient() {
	user := "Client"
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.1.68:37701")
	if err != nil {
		fmt.Fprintf(os.Stderr,"[%s] Error while getting UDP Address: %+v\n", user, err)
		os.Exit(-1)
	} else {
		fmt.Printf("Successfully resolved server address: %q\n", serverAddr)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr,"[%s] Error while dialing server: %+v\n", user, err)
		os.Exit(-3)
	} else {
		fmt.Printf("Successfully dialed connection to server. Connection: [remote: %q, local: %q]\n", conn.RemoteAddr(), conn.LocalAddr())
	}
	defer conn.Close()

	for i:=0; i<10000; i++ {
		msg := "Message number " + strconv.Itoa(i)
		buf := []byte(msg)
		fmt.Printf("Sending message: [%s]\n", msg)
		_,errs := conn.Write(buf)
		if errs != nil {
			fmt.Fprintf(os.Stderr,"Error while writing: %+v\n", errs)
		} else {
			fmt.Printf("Send packets to server.\n")
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func StartUdpServer() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":37701")
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while getting UDP Address: %+v\n", err)
	} else {
		fmt.Printf("Successfully resolved server address: %q\n", serverAddr)
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while listening at UDP Address %q: %+v\n", serverAddr, err)
	} else {
		fmt.Printf("Successfully started listening at address: %q\n", serverAddr)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		fmt.Println("Awaiting packets...")
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error while reading UDP buffer: %+v\n", err)
		} else {
			fmt.Printf("[%s] %s\n", addr.IP.String(), string(buffer[0:n]))
		}
	}
}