package utils

import "unsafe"

func BigEndian() (bool) {
	var i int = 0x1
	bs := (*[int(unsafe.Sizeof(0))]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return true
	} else {
		return false
	}
}
