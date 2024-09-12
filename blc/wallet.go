package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

const (
	// 密钥长度
	version            = byte(0x00)
	addressCheckSumLen = 4
)

type Wallet struct {
	//私钥 椭圆算法
	PrivateKey ecdsa.PrivateKey
	//公钥
	PublicKey []byte
}

// 创建 钱包
func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()

	fmt.Println("私钥：", privateKey)
	fmt.Println("公钥：", publicKey)

	return &Wallet{privateKey, publicKey}
}

// 通过私钥生成公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		return ecdsa.PrivateKey{}, nil
	}
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publicKey
}

// 获取地址
func (w *Wallet) GetAddress() string {
	//1 . hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	// 4个字节 + 20个字节 = 24个字节
	// version {0} + hash160 +4个字节 -> 25个字节
	address := append([]byte{version}, ripemd160Hash...)
	//公钥hash 2次sha256
	sha256Hash := Sha256Hash(Sha256Hash(address))
	// 取前4个字节
	checkSum := sha256Hash[:addressCheckSumLen]
	// 25个字节 + 4个字节 = 29个字节
	address = append(address, checkSum...)
	return string(Base58Encode(address))
}

func Ripemd160Hash(publickey []byte) []byte {
	// 1. sha256
	sha256Hash := Sha256Hash(publickey)
	// 2. ripemd160
	return Ripemd160HashUtils(sha256Hash)
}

// 检查钱包的地址是否合法
func IsValidAddress(address string) bool {
	// 1. base58解码
	decodeBytes := Base58Decode([]byte(address))
	fmt.Printf("decodeBytes:%x\n", decodeBytes)
	// 2. 取前21个字节
	pubKeyHash := decodeBytes[:len(decodeBytes)-addressCheckSumLen]
	// 3. 取后4个字节	作为校验码
	checkSum := decodeBytes[len(decodeBytes)-addressCheckSumLen:]
	// 4. 取前20个字节sha256
	sha256Hash := Sha256Hash(Sha256Hash(pubKeyHash))
	// 5. 取前4个字节
	checkSum2 := sha256Hash[:addressCheckSumLen]
	// 6. 比较校验码
	return bytes.Compare(checkSum, checkSum2) == 0
}
