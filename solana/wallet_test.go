package solana

import (
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/icodeface/hdkeyring"
	"math/big"
	"testing"
)

func TestNewWallet(t *testing.T) {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := NewWallet(mnemonic)
	if err != nil {
		t.Error(err)
		return
	}
	solPath, _ := hdkeyring.ParseDerivationPath("m/44'/501'/0'/0'")
	account, err := wallet.DeriveAccount(solPath)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(account.Address)

	client := NewClient(rpc.DevNet_RPC)

	// Airdrop 5 SOL to the new account:
	//out, err := client.RequestAirdrop(
	//	context.TODO(),
	//	account.PrivateKey.PublicKey(),
	//	solana.LAMPORTS_PER_SOL*1,
	//	rpc.CommitmentFinalized,
	//)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("airdrop transaction signature:", out)

	accountTo := "6hZqw492xow22UqCRW7NUZJzoPRzBTUdM2fqHN2oy76a"
	amount := big.NewInt(123000000)

	// Send transaction, and wait for confirmation:
	sig, err := account.Transfer(client, accountTo, amount)
	if err != nil {
		panic(err)
	}
	fmt.Println(sig)
}
