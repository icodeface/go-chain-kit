package filecoin

import (
	"github.com/filecoin-project/go-address"
	_ "github.com/icodeface/chain-kit/filecoin/sigs/secp" // enable secp signatures
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

	addr, err := address.NewSecp256k1Address(hdkeyring.ECDSAPublicKeyBytes(&pk.PublicKey))
	if err != nil {
		return nil, err
	}
	return &Account{
		PrivateKey: pk,
		Address:    addr,
	}, nil
}
