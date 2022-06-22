package solana

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://solana-mainnet.phantom.tech")

	height, err := client.GetBlockHeight(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(height)
}
