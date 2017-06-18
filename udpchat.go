package main

import (
	"net"
	"fmt"
	"strconv"
	"os"
	"github.com/fire00f1y/udpchat/message"
	"strings"
	"bufio"
	"time"
)

func main() {
	go StartUdpServer()
	time.Sleep(1*time.Second)
	RunUdpClient()
}

func RunUdpClient() {
	var user string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	u, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Shit went wrong: %+v\n", err)
		user = "Unknown user"
	} else  {
		user = strings.TrimRight(u, "\n")
	}
	serverAddr, err := net.ResolveUDPAddr("udp", "99.141.152.69:37701")
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
		fmt.Printf("[Client] Successfully dialed connection to server. Connection: [remote: %q, local: %q]\n",
			conn.RemoteAddr(), conn.LocalAddr())
	}
	defer conn.Close()
	index := strings.LastIndex(conn.LocalAddr().String(), ":")
	port := conn.LocalAddr().String()[index+1:]

	fmt.Println("Chat started. Type \"exit()\" to quit the application:")
	fmt.Println("====================================")
	for {
		packet := message.Packet{}
		packet.User = user
		i, err := strconv.Atoi(port)
		if err!= nil {
			packet.Port = 37701
		} else {
			packet.Port = uint16(i)
		}
		fmt.Printf("[%s] $ ", user)
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading commandline input: %+v\n", err)
			continue
		} else {
			s = strings.TrimRight(s, "\n")
			if s == "exit()" {
				fmt.Println("Exiting application...")
				os.Exit(0)
			}
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
		time.Sleep(500*time.Millisecond)
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

	fmt.Println("======= Now receiving messages =======")
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr,"[Server] Error while reading UDP buffer: %+v\n", err)
		} else {
			packet := message.New(buffer[0:n])
			fmt.Printf("Message from %s:\n%s\n", addr.String(), packet.Print())
		}
	}
}