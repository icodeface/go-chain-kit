package filecoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/icodeface/chain-kit/filecoin/types"
	"github.com/ipfs/go-cid"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

type clientRequest struct {
	Id      int64         `json:"id"`
	Version string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (r *clientRequest) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

type clientResponse struct {
	Id      uint64           `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   interface{}      `json:"error,omitempty"`
}

func (c *clientResponse) ReadFromResult(x interface{}) error {
	if x == nil {
		return nil
	}
	return json.Unmarshal(*c.Result, x)
}

type Client struct {
	addr  string
	token string
	id    int64
}

func NewClient(addr string, token string) *Client {
	return &Client{
		addr:  addr,
		token: token,
	}
}

// SetToken set Authorization token
func (c *Client) SetToken(token string) *Client {
	c.token = token
	return c
}

// Namespace Filecoin
func (c *Client) FilecoinMethod(method string) string {
	return fmt.Sprintf("Filecoin.%s", method)
}

// Request call RPC method
func (c *Client) Request(ctx context.Context, method string, result interface{}, params ...interface{}) error {
	request := &clientRequest{
		Id:      atomic.AddInt64(&c.id, 1),
		Version: "2.0",
		Method:  method,
		Params:  params,
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.addr, bytes.NewReader(request.Bytes()))
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	response := &clientResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return err
	}
	if response.Error != nil {
		return fmt.Errorf("jsonrpc call: %v", response.Error)
	}
	if response.Result == nil {
		return nil
	}

	return response.ReadFromResult(result)
}

// ChainGetMessage reads a message referenced by the specified CID from the chain blockstore.
func (c *Client) ChainGetMessage(ctx context.Context, id cid.Cid) (*types.Message, error) {
	var message *types.Message
	return message, c.Request(ctx, c.FilecoinMethod("ChainGetMessage"), &message, id)
}

// ChainGetBlockMessages returns messages stored in the specified block.
func (c *Client) ChainGetBlockMessages(ctx context.Context, id cid.Cid) (*types.BlockMessages, error) {
	var bm *types.BlockMessages
	return bm, c.Request(ctx, c.FilecoinMethod("ChainGetBlockMessages"), &bm, id)
}

// ChainHead returns the current head of the chain.
func (c *Client) ChainHead(ctx context.Context) (*types.TipSet, error) {
	var ts *types.TipSet
	return ts, c.Request(ctx, c.FilecoinMethod("ChainHead"), &ts)
}

// ChainGetTipSetByHeight looks back for a tipset at the specified epoch. If there are no blocks at the specified epoch, a tipset at an earlier epoch will be returned.
func (c *Client) ChainGetTipSetByHeight(ctx context.Context, height int64, tsk types.TipSetKey) (*types.TipSet, error) {
	var ts *types.TipSet
	return ts, c.Request(ctx, c.FilecoinMethod("ChainGetTipSetByHeight"), &ts, height, tsk)
}

// ChainExport returns a stream of bytes with CAR dump of chain data.
func (c *Client) ChainExport(ctx context.Context, tsk types.TipSetKey) ([]byte, error) {
	var result []byte
	return result, c.Request(ctx, c.FilecoinMethod("ChainExport"), &result, tsk)
}

// ChainGetBlock returns the block specified by the given CID.
func (c *Client) ChainGetBlock(ctx context.Context, id cid.Cid) (*types.BlockHeader, error) {
	var bh *types.BlockHeader
	return bh, c.Request(ctx, c.FilecoinMethod("ChainGetBlock"), &bh, id)
}

// ChainGetGenesis returns the genesis tipset.
func (c *Client) ChainGetGenesis(ctx context.Context) (*types.TipSet, error) {
	var ts *types.TipSet
	return ts, c.Request(ctx, c.FilecoinMethod("ChainGetGenesis"), &ts)
}

// ChainGetNode
func (c *Client) ChainGetNode(ctx context.Context, p string) (*types.IpldObject, error) {
	var ipld *types.IpldObject
	return ipld, c.Request(ctx, c.FilecoinMethod("ChainGetNode"), &ipld, p)
}

// ChainGetParentMessages returns messages stored in parent tipset of the specified block.
func (c *Client) ChainGetParentMessages(ctx context.Context, id cid.Cid) ([]types.ParentMessage, error) {
	var msgs []types.ParentMessage
	return msgs, c.Request(ctx, c.FilecoinMethod("ChainGetParentMessages"), &msgs, id)
}

// ChainGetParentReceipts returns receipts for messages in parent tipset of the specified block.
func (c *Client) ChainGetParentReceipts(ctx context.Context, id cid.Cid) ([]*types.MessageReceipt, error) {
	var mrs []*types.MessageReceipt
	return mrs, c.Request(ctx, c.FilecoinMethod("ChainGetParentReceipts"), &mrs, id)
}

// ChainGetPath returns a set of revert/apply operations needed to get from one tipset to another
func (c *Client) ChainGetPath(ctx context.Context, from types.TipSetKey, to types.TipSetKey) (*types.HeadChange, error) {
	var hc *types.HeadChange
	return hc, c.Request(ctx, c.FilecoinMethod("ChainGetPath"), &hc, from, to)
}

// ChainGetRandomnessFromBeacon is used to sample the beacon for randomness.
func (c *Client) ChainGetRandomnessFromBeacon(ctx context.Context, tsk types.TipSetKey, personalization int64, randEpoch int64, entropy []byte) ([]byte, error) {
	var result []byte
	return result, c.Request(ctx, c.FilecoinMethod("ChainGetRandomnessFromBeacon"), &result, tsk, personalization, randEpoch, entropy)
}

// ChainGetRandomnessFromTickets is used to sample the chain for randomness.
func (c *Client) ChainGetRandomnessFromTickets(ctx context.Context, tsk types.TipSetKey, personalization int64, randEpoch int64, entropy []byte) ([]byte, error) {
	var result []byte
	return result, c.Request(ctx, c.FilecoinMethod("ChainGetRandomnessFromTickets"), &result, tsk, personalization, randEpoch, entropy)
}

// ChainGetTipSet returns the tipset specified by the given TipSetKey.
func (c *Client) ChainGetTipSet(ctx context.Context, tsk types.TipSetKey) (*types.TipSet, error) {
	var ts *types.TipSet
	return ts, c.Request(ctx, c.FilecoinMethod("ChainGetTipSet"), &ts, tsk)
}

// ChainHasObj checks if a given CID exists in the chain blockstore.
func (c *Client) ChainHasObj(ctx context.Context, o cid.Cid) (bool, error) {
	var ok bool
	return ok, c.Request(ctx, c.FilecoinMethod("ChainHasObj"), &ok, o)
}

// ChainReadObj reads ipld nodes referenced by the specified CID from chain blockstore and returns raw bytes.
func (c *Client) ChainReadObj(ctx context.Context, obj cid.Cid) ([]byte, error) {
	var result []byte
	return result, c.Request(ctx, c.FilecoinMethod("ChainReadObj"), &result, obj)
}

// ChainSetHead forcefully sets current chain head. Use with caution.
func (c *Client) ChainSetHead(ctx context.Context, tsk types.TipSetKey) error {
	return c.Request(ctx, c.FilecoinMethod("ChainSetHead"), nil, tsk)
}

// ChainStatObj returns statistics about the graph referenced by 'obj'. If 'base' is also specified, then the returned stat will be a diff between the two objects.
func (c *Client) ChainStatObj(ctx context.Context, obj, base cid.Cid) (types.ObjStat, error) {
	var os types.ObjStat
	return os, c.Request(ctx, c.FilecoinMethod("ChainStatObj"), &os, obj, base)
}

// ChainTipSetWeight computes weight for the specified tipset.
func (c *Client) ChainTipSetWeight(ctx context.Context, tsk types.TipSetKey) (decimal.Decimal, error) {
	var d decimal.Decimal
	return d, c.Request(ctx, c.FilecoinMethod("ChainTipSetWeight"), &d, tsk)
}

// BeaconGetEntry returns the beacon entry for the given filecoin epoch. If the entry has not yet been produced, the call will block until the entry becomes available
func (c *Client) BeaconGetEntry(ctx context.Context, epoch int64) (*types.BeaconEntry, error) {
	var be *types.BeaconEntry
	return be, c.Request(ctx, c.FilecoinMethod("BeaconGetEntry"), &be, epoch)
}

// GasEstimateGasLimit estimates gas used by the message and returns it. It fails if message fails to execute.
func (c *Client) GasEstimateGasLimit(ctx context.Context, message *types.Message, cids []*cid.Cid) (int64, error) {
	var gasLimit int64
	return gasLimit, c.Request(ctx, c.FilecoinMethod("GasEstimateGasLimit"), &gasLimit, message, cids)
}

// GasEstimateMessageGas estimates gas values for unset message gas fields
func (c *Client) GasEstimateMessageGas(ctx context.Context, message *types.Message, spec *types.MessageSendSpec, cids []*cid.Cid) (*types.Message, error) {
	var msg *types.Message
	return msg, c.Request(ctx, c.FilecoinMethod("GasEstimateMessageGas"), &msg, message, spec, cids)
}

// MpoolPush pushes a signed message to mempool.
func (c *Client) MpoolPush(ctx context.Context, sm *types.SignedMessage) (cid.Cid, error) {
	var id cid.Cid
	return id, c.Request(ctx, c.FilecoinMethod("MpoolPush"), &id, sm)
}

// MpoolGetNonce 获取指定发送账号的下一个nonce值
func (c *Client) MpoolGetNonce(ctx context.Context, address address.Address) (nonce uint64, err error) {
	return nonce, c.Request(ctx, c.FilecoinMethod("MpoolGetNonce"), &nonce, address)
}

// StateGetActor returns the indicated actor's nonce and balance.
func (c *Client) StateGetActor(ctx context.Context, addr address.Address, cids []*cid.Cid) (*types.Actor, error) {
	var actor *types.Actor
	return actor, c.Request(ctx, c.FilecoinMethod("StateGetActor"), &actor, addr, cids)
}

// StateGetReceipt returns the message receipt for the given message
func (c *Client) StateGetReceipt(ctx context.Context, id cid.Cid, cids []*cid.Cid) (*types.MessageReceipt, error) {
	var mr *types.MessageReceipt
	return mr, c.Request(ctx, c.FilecoinMethod("StateGetReceipt"), &mr, id, cids)
}

// StateReplay returns the result of executing the indicated message, assuming it was executed in the indicated tipset.
func (c *Client) StateReplay(ctx context.Context, tsk types.TipSetKey, mc cid.Cid) (*types.InvocResult, error) {
	var result *types.InvocResult
	return result, c.Request(ctx, c.FilecoinMethod("StateReplay"), &result, tsk, mc)
}

// StateSearchMsg searches for a message in the chain, and returns its receipt and the tipset where it was executed
func (c *Client) StateSearchMsg(ctx context.Context, msg cid.Cid) (*types.MsgLookup, error) {
	var msgl *types.MsgLookup
	return msgl, c.Request(ctx, c.FilecoinMethod("StateSearchMsg"), &msgl, msg)
}

// WalletBalance returns the balance of the given address at the current head of the chain.
func (c *Client) WalletBalance(ctx context.Context, addr address.Address) (abi.TokenAmount, error) {
	var balance abi.TokenAmount
	return balance, c.Request(ctx, c.FilecoinMethod("WalletBalance"), &balance, addr)
}

// WalletDefaultAddress returns the address marked as default in the wallet.
func (c *Client) WalletDefaultAddress(ctx context.Context) (address.Address, error) {
	var addr address.Address
	return addr, c.Request(ctx, c.FilecoinMethod("WalletDefaultAddress"), &addr)
}

// WalletDelete deletes an address from the wallet.
func (c *Client) WalletDelete(ctx context.Context, addr address.Address) error {
	return c.Request(ctx, c.FilecoinMethod("WalletDelete"), nil, addr)
}

// WalletExport returns the private key of an address in the wallet.
func (c *Client) WalletExport(ctx context.Context, addr address.Address) (*types.KeyInfo, error) {
	var ki *types.KeyInfo
	return ki, c.Request(ctx, c.FilecoinMethod("WalletExport"), &ki, addr)
}

// WalletHas indicates whether the given address is in the wallet.
func (c *Client) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	var has bool
	return has, c.Request(ctx, c.FilecoinMethod("WalletHas"), &has, addr)
}

// WalletImport receives a KeyInfo, which includes a private key, and imports it into the wallet.
func (c *Client) WalletImport(ctx context.Context, ki *types.KeyInfo) (address.Address, error) {
	var addr address.Address
	return addr, c.Request(ctx, c.FilecoinMethod("WalletImport"), &addr, ki)
}

// WalletList lists all the addresses in the wallet.
func (c *Client) WalletList(ctx context.Context) ([]address.Address, error) {
	var addrs []address.Address
	return addrs, c.Request(ctx, c.FilecoinMethod("WalletList"), &addrs)
}

// WalletNew creates a new address in the wallet with the given KeyType.
func (c *Client) WalletNew(ctx context.Context, typ types.KeyType) (address.Address, error) {
	var addr address.Address
	return addr, c.Request(ctx, c.FilecoinMethod("WalletNew"), &addr, typ)
}

// WalletSetDefault marks the given address as as the default one.
func (c *Client) WalletSetDefault(ctx context.Context, addr address.Address) error {
	return c.Request(ctx, c.FilecoinMethod("WalletSetDefault"), nil, addr)
}

// WalletSign signs the given bytes using the given address.
func (c *Client) WalletSign(ctx context.Context, addr address.Address, data []byte) (*crypto.Signature, error) {
	var sig *crypto.Signature
	return sig, c.Request(ctx, c.FilecoinMethod("WalletSign"), &sig, addr, data)
}

// WalletSignMessage signs the given message using the given address.
func (c *Client) WalletSignMessage(ctx context.Context, addr address.Address, message *types.Message) (*types.SignedMessage, error) {
	var sm *types.SignedMessage
	return sm, c.Request(ctx, c.FilecoinMethod("WalletSignMessage"), &sm, addr, message)
}

// WalletVerify takes an address, a signature, and some bytes, and indicates whether the signature is valid. The address does not have to be in the wallet.
func (c *Client) WalletVerify(ctx context.Context, k string, msg []byte, sig *crypto.Signature) (bool, error) {
	var ok bool
	return ok, c.Request(ctx, c.FilecoinMethod("WalletVerify"), &ok, k, msg, sig)
}

// Version Get Lotus Node version
func (c *Client) Version(ctx context.Context) (*types.Version, error) {
	var version *types.Version
	return version, c.Request(ctx, c.FilecoinMethod("Version"), &version)
}
