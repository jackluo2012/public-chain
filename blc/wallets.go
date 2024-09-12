package blc

type Wallets struct {
	Wallets map[string]*Wallet
}

// 创建钱包集合
func NewWallets() *Wallets {
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	return wallets
}

// 创建一个新的钱包
func (wallets *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()
	wallets.Wallets[address] = wallet
	return address
}
