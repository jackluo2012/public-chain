package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"math/big"

	"io/ioutil"
	"log"
	"os"
)

// 钱包信息保存到文件中
const (
	walletsFile = "./wallets.dat"
)

type Wallets struct {
	WalletsMap map[string]*Wallet `json:"WalletsMap"`
}

// SerializableWallet 解决gob: type elliptic.p256Curve has no exported fields问题
type SerializableWallet struct {
	D         *big.Int
	X, Y      *big.Int
	PublicKey []byte
}

// 检查钱包 文件是否存在
func WalletsExist() bool {
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 加载钱包文件
func (ws *Wallets) LoadWallets() error {
	if !WalletsExist() {
		return errors.New("钱包文件不存在")
	}

	fileContent, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets map[string]SerializableWallet
	//gob.Register(elliptic.P256())
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = make(map[string]*Wallet)
	//ws.Wallets = wallets.Wallets
	for k, v := range wallets {
		ws.WalletsMap[k] = &Wallet{
			PrivateKey: ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PublicKey: v.PublicKey,
		}
	}

	return nil
}

// 创建钱包集合
func NewWallets() (*Wallets, error) {
	wallets := &Wallets{}
	wallets.WalletsMap = make(map[string]*Wallet)
	err := wallets.LoadWallets()
	// 钱包信息取出来
	return wallets, err
}

// 创建一个新的钱包
func (wallets *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	wallets.WalletsMap[address] = wallet
	// 把所有钱包保存到文件
	wallets.SaveWallets()
	return address
}

// 把钱包保存到文件中
func (ws *Wallets) SaveWallets() {

	var content bytes.Buffer
	gob.Register(SerializableWallet{})

	wallets := make(map[string]SerializableWallet)
	for k, v := range ws.WalletsMap {
		wallets[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PublicKey,
		}
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletsFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
