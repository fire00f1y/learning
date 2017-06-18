package message

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"os"
	"github.com/fire00f1y/learning/utils"
	"time"
	"strings"
)

var (
	byteOrder  binary.ByteOrder
	DateFormat = "2006-01-02 15:04:05.000"
)

func init() {
	if utils.BigEndian() {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

/*
Packet byte structure:
byte  [0:0]: Username length
bytes [1:2]: Port value
bttes [3:24]: Timestamp
bytes [24:24+${Username length}]: Username
bytes [24+${Username length} + 1: ]: Message
 */
type Packet struct {
	Port      uint16	`bytes:"2"`
	User      string	`bytes:"len(.User)"`
	Timestamp string	`bytes:"23"`
	Message   string	`bytes:"len(.Message)"`
}

func (p Packet) BinaryMarshaler() ([]byte, error) {
	// Create byte arrays
	portBytes := make([]byte, 2)
	userLength := make([]byte, 1)
	timestamp := []byte(time.Now().Format(DateFormat))
	userBytes := []byte(p.User)
	mesBytes := []byte(p.Message)

	// Calculate total size and allocate output bytes
	total := len(userBytes) + len(mesBytes) + len(portBytes) + len(userLength) + len(timestamp)
	mod := 4 - (total % 4)
	returnbytes := make([]byte, total+mod)
	index := 0

	// Populate metadata bytes
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, p.Port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting port to binary: %+v\n", err)
		portBytes[0], portBytes[1] = uint8(p.Port>>8), uint8(p.Port&0xff)
	} else {
		portBytes = buffer.Bytes()
	}

	buffer = new(bytes.Buffer)
	err = binary.Write(buffer, binary.LittleEndian, uint8(len(p.User)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting user length to binary: %+v\n", err)
		userLength[0] = uint8(len(p.User) & 0xff)
	} else {
		userLength = buffer.Bytes()
	}

	// Write bytes out
	for _, v := range userLength {
		returnbytes[index] = v
		index++
	}
	for _, v := range portBytes {
		returnbytes[index] = v
		index++
	}
	for _, v := range timestamp {
		returnbytes[index] = v
		index++
	}
	for _, v := range userBytes {
		returnbytes[index] = v
		index++
	}
	for _, v := range mesBytes {
		returnbytes[index] = v
		index++
	}
	for index < total+mod {
		returnbytes[index] = ' '
		index += 1
	}

	return returnbytes, nil
}

func New(input []byte) (Packet) {
	p := Packet{}
	userLength  := input[0:1]
	var length uint8
	err := binary.Read(bytes.NewReader(userLength), binary.LittleEndian, &length)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading user name length: %+v\n", err)
	}

	portBytes := input[1:3]
	err = binary.Read(bytes.NewReader(portBytes), binary.LittleEndian, &p.Port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading in port: %+v\n", err)
	}

	timestamp := input[3:26]
	p.Timestamp = string(timestamp)

	username := input[26:(26 + length)]
	p.User = string(username)

	mes := input[26+length:]
	mesString := string(mes)
	p.Message = strings.TrimSpace(mesString)
	return p
}

func (p Packet) Print() (string) {
	return fmt.Sprintf("[%s] %s: %s\n", p.Timestamp, p.User, p.Message)
}