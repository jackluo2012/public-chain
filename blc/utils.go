package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"golang.org/x/crypto/ripemd160"
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

// json 字符串转成数组
func JsonToArr(jsonStr string) []string {
	var arr []string
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		panic(err)
	}
	return arr
}

// string 转成int64
func StrToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}

// string 转成字节 数组
func StrToBytes(str string) []byte {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return bytes
}

// 字节 数组 转成 string
func BytesToStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// Int64 转 []byte
func Int64ToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// Sha256Hash
func Sha256Hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// Ripemd160Hash
func Ripemd160HashUtils(data []byte) []byte {
	hash := ripemd160.New()
	hash.Write(data)
	return hash.Sum(nil)
}
