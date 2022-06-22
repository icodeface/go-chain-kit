package solana

import (
	"fmt"
	"github.com/icodeface/hdkeyring"
	"golang.org/x/xerrors"
)

type Wallet struct {
	keyring *hdkeyring.Keyring
}

func NewWallet(mnemonic string) (*Wallet, error) {
	keyring, err := hdkeyring.NewFromMnemonic(mnemonic, hdkeyring.KeyTypeEd25519)
	if err != nil {
		return nil, err
	}
	return &Wallet{keyring: keyring}, nil
}

func (w *Wallet) DeriveAccount(path hdkeyring.DerivationPath) (*Account, error) {
	pk, err := w.keyring.DeriveEd25519PrivateKey(path)
	if err != nil {
		return nil, xerrors.Errorf("derive private key: %w", err)
	}
	return AccountFromPrivateKey(*pk), nil
}

func DerivePath(index int64) hdkeyring.DerivationPath {
	return hdkeyring.MustParseDerivationPath(fmt.Sprintf("m/44'/501'/0'/%d'", index))
}
