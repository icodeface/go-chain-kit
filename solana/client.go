package solana

import (
	"github.com/gagliardetto/solana-go/rpc"
)

type Client = rpc.Client

func NewClient(endpoint string) *Client {
	return rpc.New(endpoint)
}
