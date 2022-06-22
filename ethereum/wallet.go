package ethereum

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icodeface/hdkeyring"
)

type Wallet struct {
	keyring *hdkeyring.Keyring
}

func NewWallet(mnemonic string) (*Wallet, error) {
	keyring, err := hdkeyring.NewFromMnemonic(mnemonic, hdkeyring.KeyTypeECDSA)
	if err != nil {
		return nil, err
	}
	return &Wallet{keyring: keyring}, nil
}

func (w *Wallet) DeriveAccount(path hdkeyring.DerivationPath) (*Account, error) {
	pk, err := w.keyring.DeriveECDSAPrivateKey(path)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(pk.PublicKey)
	return &Account{
		PrivateKey: pk,
		Address:    address,
	}, nil
}
