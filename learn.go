package main

import (
	"net"
	"fmt"
	"strconv"
	"os"
	"github.com/fire00f1y/learning/message"
	"strings"
	"bufio"
)

func main() {
	go StartUdpServer()
	RunUdpClient()
}

func RunUdpClient() {
	user := "fire00f1y"
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.1.68:37701")
	if err != nil {
		fmt.Fprintf(os.Stderr,"[Client] Error while getting UDP Address: %+v\n", err)
		os.Exit(-1)
	} else {
		fmt.Printf("[Client] Successfully resolved server address: %q\n", serverAddr)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr,"[Client] Error while dialing server: %+v\n", err)
		os.Exit(-3)
	} else {
		fmt.Printf("[Client] Successfully dialed connection to server. Connection: [remote: %q, local: %q]\n", conn.RemoteAddr(), conn.LocalAddr())
	}
	index := strings.LastIndex(conn.LocalAddr().String(), ":")
	port := conn.LocalAddr().String()[index+1:]
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for i:=0; i<10000; i++ {
		fmt.Print(".")
		packet := message.Packet{}
		packet.User = user
		i, err = strconv.Atoi(port)
		if err!= nil {
			packet.Port = 37701
		} else {
			packet.Port = uint16(i)
		}
		fmt.Print("$ ")
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading commandline input: %+v\n", err)
			continue
		} else {
			packet.Message = s
		}
		buf, err := packet.BinaryMarshaler()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while marshalling message: %+v\n", err)
			continue
		}
		_,errs := conn.Write(buf)
		if errs != nil {
			fmt.Fprintf(os.Stderr,"[Client] Error while writing: %+v\n", errs)
		}
	}
}

func StartUdpServer() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":37701")
	if err != nil {
		fmt.Fprintf(os.Stderr,"[Server] Error while getting UDP Address: %+v\n", err)
	} else {
		fmt.Printf("[Server] Successfully resolved server address: %q\n", serverAddr)
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr,"[Server] Error while listening at UDP Address %q: %+v\n", serverAddr, err)
	} else {
		fmt.Printf("[Server] Successfully started listening at address: %q\n", serverAddr)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		fmt.Println("Awaiting packets...")
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr,"[Server] Error while reading UDP buffer: %+v\n", err)
		} else {
			packet := message.New(buffer[0:n])
			fmt.Printf("Message from %q:\n%s\n", addr, packet.Print())
		}
	}
}