package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type RpcRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      uint64 `json:"id"`
}

type RpcResponse struct {
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	ID      uint64          `json:"id"`
	Error   *RpcError       `json:"error"`
}

type RpcError struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
	Data    any    `json:"data"`
}

func (e *RpcError) Error() string {
	return fmt.Sprintf("rpc error status=%d message=%s", e.Code, e.Message)
}

type RpcClient struct {
	url     string
	client  *http.Client
	chainID string
}

func NewRpcClient(ctx context.Context, url string) (*RpcClient, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	rpc := &RpcClient{
		url:    url,
		client: client,
	}
	chainID, err := rpc.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	rpc.chainID = chainID
	return rpc, nil
}

func (c *RpcClient) call(ctx context.Context, method string, params ...any) (*RpcResponse, error) {
	body, err := json.Marshal(&RpcRequest{
		ID:      0,
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var rpcRes RpcResponse
	if err := json.Unmarshal(resData, &rpcRes); err != nil {
		return nil, err
	}
	return &rpcRes, nil
}

// BlockNumber returns the number of most recent block.
func (c *RpcClient) BlockNumber(ctx context.Context) (uint64, error) {
	res, err := c.call(ctx, "eth_blockNumber")
	if err != nil {
		return 0, err
	}
	if res.Error != nil {
		return 0, res.Error
	}
	var number string
	if err := json.Unmarshal(res.Result, &number); err != nil {
		return 0, err
	}
	return strconv.ParseUint(number, 0, 64)
}

// GetBlockByNumber returns information about a block by block number.
func (c *RpcClient) GetBlockByNumber(ctx context.Context, number uint64, fullTx bool) (map[string]any, error) {
	res, err := c.call(ctx, "eth_getBlockByNumber", fmt.Sprintf("0x%x", number), fullTx)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	var block map[string]any
	if err := json.Unmarshal(res.Result, &block); err != nil {
		return nil, err
	}
	return block, nil
}

// ChainID returns the unique identifier for the blockchain.
func (c *RpcClient) ChainID(ctx context.Context) (string, error) {
	if c.chainID != "" {
		return c.chainID, nil
	}
	res, err := c.call(ctx, "eth_chainId")
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", res.Error
	}
	var chainID string
	if err := json.Unmarshal(res.Result, &chainID); err != nil {
		return "", err
	}
	return chainID, nil
}
