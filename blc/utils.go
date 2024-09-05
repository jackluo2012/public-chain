package blc

import (
	"bytes"
	"encoding/binary"
)

func IntToHex(num int64) []byte {
	// 1. 将int64转换为[]byte
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}
