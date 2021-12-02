package ethereum_test

import (
	"context"
	"fmt"
	"github.com/icodeface/chain-kit/ethereum"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := ethereum.NewClient(context.Background(), "https://mainnet.infura.io/v3/6b851498cebb4893b436e986cf1f5458")
	if err != nil {
		t.Error(err)
		return
	}
	head, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(head)
}
