package filecoin

import (
	"context"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/icodeface/chain-kit/filecoin/sigs"
	_ "github.com/icodeface/chain-kit/filecoin/sigs/secp" // enable secp signatures
	"github.com/icodeface/chain-kit/filecoin/types"
	"github.com/icodeface/hdwallet"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
	"math/big"
	"time"
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

func (w *Wallet) DeriveAddress(path hdwallet.DerivationPath) (address.Address, error) {
	publicKeyECDSA, err := w.hd.DerivePublicKey(path)
	if err != nil {
		return address.Address{}, err
	}

	addr, err := address.NewSecp256k1Address(hdwallet.PublicKeyBytes(publicKeyECDSA))
	if err != nil {
		return address.Address{}, err
	}
	return addr, nil
}

func (w *Wallet) MustDeriveAddress(path hdwallet.DerivationPath) address.Address {
	addr, err := w.DeriveAddress(path)
	if err != nil {
		panic(err)
	}
	return addr
}

func (w *Wallet) SignMessage(path hdwallet.DerivationPath, msg *types.Message) (*types.SignedMessage, error) {
	mb, err := msg.ToStorageBlock()
	if err != nil {
		return nil, xerrors.Errorf("serializing message: %w", err)
	}

	pk, err := w.hd.DerivePrivateKey(path)
	if err != nil {
		return nil, xerrors.Errorf("derive private key: %w", err)
	}

	sig, err := sigs.Sign(crypto.SigTypeSecp256k1, hdwallet.PrivateKeyBytes(pk), mb.Cid().Bytes())
	if err != nil {
		return nil, fmt.Errorf("sign message: %w", err)
	}

	return &types.SignedMessage{
		Message:   msg,
		Signature: sig,
	}, nil
}

func (w *Wallet) Transfer(client *Client, fromPath hdwallet.DerivationPath, to string, amount *big.Int) (cid.Cid, error) {
	toAddr, err := address.NewFromString(to)
	if err != nil {
		return cid.Undef, err
	}
	if toAddr == address.Undef {
		return cid.Undef, errors.New("empty address")
	}
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return cid.Undef, errors.New("invalid value")
	}

	fromAddr := w.MustDeriveAddress(fromPath)

	msg := &types.Message{
		From:   fromAddr,
		To:     toAddr,
		Value:  BigIntToTokenAmount(amount),
		Method: types.MethodSend,
		Params: nil,
	}

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err = client.GasEstimateMessageGas(timeout, msg, nil, nil)
	if err != nil {
		return cid.Undef, xerrors.Errorf("GasEstimateMessageGas error: %w", err)
	}
	if msg.GasPremium.GreaterThan(msg.GasFeeCap) {
		return cid.Undef, xerrors.Errorf("After estimation, GasPremium is greater than GasFeeCap")
	}

	timeout, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	b, err := client.WalletBalance(timeout, msg.From)
	if err != nil {
		return cid.Undef, xerrors.Errorf("getting origin balance: %w", err)
	}
	if b.LessThan(msg.Value) {
		return cid.Undef, xerrors.Errorf("not enough funds: %s < %s", b, msg.Value)
	}

	timeout, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nonce, err := client.MpoolGetNonce(timeout, fromAddr)
	if err != nil {
		return cid.Undef, xerrors.Errorf("mpool get nonce: %w", err)
	}
	msg.Nonce = nonce

	signed, err := w.SignMessage(DerivePath(0, 0), msg)
	if err != nil {
		return cid.Undef, err
	}

	timeout, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	hash, err := client.MpoolPush(timeout, signed)
	if err != nil {
		return cid.Undef, xerrors.Errorf("failed to push message: %w", err)
	}
	return hash, nil
}
