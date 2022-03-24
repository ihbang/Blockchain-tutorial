package main

const walletFile = "wallet.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()

	ws.Wallets[address] = wallet

	return address
}

func (ws Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}
