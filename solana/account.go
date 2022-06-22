package solana

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"math/big"
)

type Account struct {
	PrivateKey solana.PrivateKey
	Address    string
}

func AccountFromPrivateKey(rawKey []byte) *Account {
	priv := solana.PrivateKey(rawKey)
	pub := priv.PublicKey()
	address := pub.String()
	return &Account{
		PrivateKey: priv,
		Address:    address,
	}
}

func (account *Account) PublicKey() solana.PublicKey {
	return account.PrivateKey.PublicKey()
}

func (account *Account) Transfer(rpcClient *Client, to string, amount *big.Int) (*solana.Signature, error) {
	accountTo, err := solana.PublicKeyFromBase58(to)
	if err != nil {
		return nil, err
	}

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount.Uint64(),
				account.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(account.PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if account.PublicKey().Equals(key) {
				return &account.PrivateKey
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to sign transaction: %w", err)
	}

	// Send transaction
	sig, err := rpcClient.SendTransactionWithOpts(context.TODO(), tx, false, rpc.CommitmentFinalized)
	return &sig, err
}
