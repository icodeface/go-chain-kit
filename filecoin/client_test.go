package filecoin

import (
	"context"
	"github.com/filecoin-project/go-address"
	"testing"
)

// The Lotus Node
// The default token is in ~/.lotus/token
func testClient() *Client {
	return NewClient("https://fil-rpc.foxnb.net/rpc/v0", "")
}

// 测试RpcClient
func TestClient_Request(t *testing.T) {
	c := NewClient("https://eth-mainnet.token.im", "")
	var blockNumber string
	if err := c.Request(context.Background(), "eth_blockNumber", &blockNumber); err != nil {
		t.Error(err)
	}

	t.Log(blockNumber)

	var tr struct {
		BlockHash   string `json:"blockHash"`
		BlockNumber string `json:"blockNumber"`
	}
	if err := c.Request(context.Background(), "eth_getTransactionReceipt", &tr, "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0"); err != nil {
		t.Error(err)
	}

	t.Log(tr.BlockHash)
	t.Log(tr.BlockNumber)
}

// 根据消息Cid获取消息
func TestClient_ChainGetMessage(t *testing.T) {
	c := testClient()

	head, err := c.ChainHead(context.Background())
	if err != nil {
		t.Error(err)
	}

	ts, err := c.ChainGetTipSetByHeight(context.Background(), head.Height-1, nil)
	if err != nil {
		t.Error(err)
	}
	bm, err := c.ChainGetBlockMessages(context.Background(), ts.Cids[0])

	msg, err := c.ChainGetMessage(context.Background(), bm.Cids[0])
	if err != nil {
		t.Error(err)
	}

	t.Log(msg)
	t.Log(msg.Cid().String())
}

// 获取当前头部高度
func TestClient_ChainHead(t *testing.T) {
	c := testClient()

	ts, err := c.ChainHead(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Log(ts.Height)

	for _, n := range ts.Cids {
		bm, err := c.ChainGetBlockMessages(context.Background(), n)
		if err != nil {
			t.Error(err)
		}
		for index, msg := range bm.BlsMessages {
			t.Log(bm.Cids[index], msg)
		}
	}
}

// 根据高度遍历区块所有交易
func TestClient_ChainGetTipSetByHeight(t *testing.T) {
	c := testClient()

	head, err := c.ChainHead(context.Background())
	if err != nil {
		t.Error(err)
	}

	ts, err := c.ChainGetTipSetByHeight(context.Background(), head.Height-1, nil)
	if err != nil {
		t.Error(err)
	}
	for _, n := range ts.Cids {
		bm, err := c.ChainGetBlockMessages(context.Background(), n)
		if err != nil {
			t.Error(err)
		}
		for index, msg := range bm.BlsMessages {
			t.Log(bm.Cids[index], msg)
		}
	}
}

// 遍历区块的 parentMessages
func TestClient_ChainGetParentMessages(t *testing.T) {
	c := testClient()

	head, err := c.ChainHead(context.Background())
	if err != nil {
		t.Error(err)
	}
	blockHeight := head.Height - 2

	tipSet, err := c.ChainGetTipSetByHeight(context.Background(), blockHeight, nil)
	if err != nil {
		t.Error(err)
	}
	//同一个 tipSet 下的 block 的 parentMessages 相同
	pms, err := c.ChainGetParentMessages(context.Background(), tipSet.Cids[0])
	if err != nil {
		t.Error(err)
	}
	t.Log(len(pms))
	for _, pm := range pms {
		address.CurrentNetwork = address.Mainnet
		from := pm.Message.From.String()
		to := pm.Message.To.String()
		value := pm.Message.Value
		t.Log(pm.Cid.String())
		t.Log(from)
		t.Log(to)
		t.Log(ToFil(value).String())
	}
}

// 查询消息/交易执行状态
func TestClient_StateGetReceipt(t *testing.T) {
	c := testClient()

	head, err := c.ChainHead(context.Background())
	if err != nil {
		t.Error(err)
	}

	ts, err := c.ChainGetTipSetByHeight(context.Background(), head.Height-2, nil)
	if err != nil {
		t.Error(err)
	}

	bm, err := c.ChainGetBlockMessages(context.Background(), ts.Cids[0])
	if err != nil {
		t.Error(err)
	}

	mr, err := c.StateGetReceipt(context.Background(), bm.Cids[0], nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(mr)
}

// 查询消息状态
// Receipt 为空表示未执行
func TestClient_StateSearchMsg(t *testing.T) {
	c := testClient()

	head, err := c.ChainHead(context.Background())
	if err != nil {
		t.Error(err)
	}

	ts, err := c.ChainGetTipSetByHeight(context.Background(), head.Height-2, nil)
	if err != nil {
		t.Error(err)
	}

	bm, err := c.ChainGetBlockMessages(context.Background(), ts.Cids[0])
	if err != nil {
		t.Error(err)
	}

	msg, err := c.StateSearchMsg(context.Background(), bm.Cids[0])
	if err != nil {
		t.Error(err)
	}

	if msg == nil {
		t.Log("nil")
	} else {
		t.Log(msg)
	}
}

func TestClient_StateGetActor(t *testing.T) {
	c := testClient()

	address.CurrentNetwork = address.Mainnet

	addr, _ := address.NewFromString("f3qx3jo74v6d6z35qhfeax3xozsegzliowrrchuyumshnwb2kz66xajhl55pxjr5xvvpeggioytv7uko5hpzga")

	actor, err := c.StateGetActor(context.Background(), addr, nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(actor.Nonce)

	nonce, err := c.MpoolGetNonce(context.Background(), addr)
	if err != nil {
		t.Error(err)
	}

	t.Log(nonce)
}

// 查询钱包余额
func TestClient_WalletBalance(t *testing.T) {
	c := testClient()

	addr, _ := address.NewFromString("f1ntod647g54mv7pqbkniqnyov6k7thr2uxdec42i")
	b, err := c.WalletBalance(context.Background(), addr)
	if err != nil {
		t.Error(err)
	}

	t.Log(ToFil(b))
}
