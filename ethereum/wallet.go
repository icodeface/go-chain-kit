package ethereum

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icodeface/hdwallet"
	"golang.org/x/xerrors"
	"math/big"
)

type Wallet struct {
	hd *hdwallet.Wallet
}

func NewWallet(mnemonic string) (*Wallet, error) {
	hd, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return &Wallet{hd: hd}, nil
}

func (w *Wallet) DeriveAddress(path hdwallet.DerivationPath) (common.Address, error) {
	publicKeyECDSA, err := w.hd.DerivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func (w *Wallet) MustDeriveAddress(path hdwallet.DerivationPath) common.Address {
	address, err := w.DeriveAddress(path)
	if err != nil {
		panic(err)
	}
	return address
}

func (w *Wallet) SignTransaction(path hdwallet.DerivationPath, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	pk, err := w.hd.DerivePrivateKey(path)
	if err != nil {
		return nil, xerrors.Errorf("derive private key: %w", err)
	}

	signer := types.NewEIP155Signer(chainID)
	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, signer, pk)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (w *Wallet) NewKeyedTransactorWithChainID(path hdwallet.DerivationPath, chainID *big.Int) (*bind.TransactOpts, error) {
	pk, err := w.hd.DerivePrivateKey(path)
	if err != nil {
		return nil, xerrors.Errorf("derive private key: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, xerrors.Errorf("derive private key: %w", err)
	}
	auth.GasLimit = uint64(200000)
	auth.Value = big.NewInt(0)
	auth.Nonce = nil
	return auth, nil
}
