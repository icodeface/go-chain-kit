package filecoin

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/icodeface/chain-kit/filecoin/sigs"
	"github.com/icodeface/chain-kit/filecoin/types"
	"github.com/icodeface/hdkeyring"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
	"math/big"
	"time"
)

type Account struct {
	PrivateKey *ecdsa.PrivateKey
	Address    address.Address
}

func (account *Account) SignMessage(msg *types.Message) (*types.SignedMessage, error) {
	mb, err := msg.ToStorageBlock()
	if err != nil {
		return nil, xerrors.Errorf("serializing message: %w", err)
	}

	sig, err := sigs.Sign(crypto.SigTypeSecp256k1, hdkeyring.ECDSAPrivateKeyBytes(account.PrivateKey), mb.Cid().Bytes())
	if err != nil {
		return nil, fmt.Errorf("sign message: %w", err)
	}

	return &types.SignedMessage{
		Message:   msg,
		Signature: sig,
	}, nil
}

func (account *Account) Transfer(client *Client, to string, amount *big.Int) (cid.Cid, error) {
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

	msg := &types.Message{
		From:   account.Address,
		To:     toAddr,
		Value:  BigIntToTokenAmount(amount),
		Method: types.MethodSend,
		Params: nil,
	}

	timeout, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	msg, err = client.GasEstimateMessageGas(timeout, msg, nil, nil)
	if err != nil {
		return cid.Undef, xerrors.Errorf("GasEstimateMessageGas error: %w", err)
	}
	if msg.GasPremium.GreaterThan(msg.GasFeeCap) {
		return cid.Undef, xerrors.Errorf("After estimation, GasPremium is greater than GasFeeCap")
	}

	timeout, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	b, err := client.WalletBalance(timeout, msg.From)
	if err != nil {
		return cid.Undef, xerrors.Errorf("getting origin balance: %w", err)
	}
	if b.LessThan(msg.Value) {
		return cid.Undef, xerrors.Errorf("not enough funds: %s < %s", b, msg.Value)
	}

	timeout, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	nonce, err := client.MpoolGetNonce(timeout, account.Address)
	if err != nil {
		return cid.Undef, xerrors.Errorf("mpool get nonce: %w", err)
	}
	msg.Nonce = nonce

	signed, err := account.SignMessage(msg)
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
