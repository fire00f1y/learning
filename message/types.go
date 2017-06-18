package message

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"os"
	"github.com/fire00f1y/learning/utils"
)

var byteOrder binary.ByteOrder

func init() {
	if utils.BigEndian() {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

type Packet struct {
	User string
	Message string
	Port uint16
}

/*
Packet byte structure:
bytes [0:1]: Port value
byte  [2:2]: Username length
bytes [3:4]: Message length
bytes [5:5+${Username length}]: Username
bytes [5+${Username length}: ]: Message
 */
func (p Packet) Stream() ([]byte) {
	bytes := make([]byte, 0)

	userBytes := []byte(p.User)
	mesBytes := []byte(p.Message)

	// Serialize port
	buffer := new(bytes.Buffer)
	portBytes := make([]byte, 2)
	err := binary.Write(buffer, byteOrder, p.Port)
	if err != nil {
		a, b := uint8(p.Port>>8), uint8(p.Port&0xff)
		bytes = append(bytes, a, b)
		fmt.Fprintf(os.Stderr, "Error converting port %d to LittleEndian byte array: %+v\nConverted via bitshifting: %x\n", p.Port, err, portBytes)
	} else {
		portBytes = buffer.Bytes()
		for _, v := range portBytes {
			bytes = append(bytes, v)
		}
	}

	buffer = new(bytes.Buffer)
	userSizeBytes := make([]byte, 1)
	err = binary.Write(buffer, byteOrder, len(p.User))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing byte array of username: %+v\n", err)
		userSizeBytes = make([]byte, 256)
		for i, v := range userBytes {

		}
	}

	return bytes
}

func (p *Packet) Read(bytes []byte) {

}