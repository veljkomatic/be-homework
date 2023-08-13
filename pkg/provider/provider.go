package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/veljkomatic/be-homework/common/jsonrpc"
)

const (
	httpTimeout = 10 * time.Second
)

// Provider is responsible for providing blockchain data
type Provider interface {
	// GetLatestBlockNumber returns the latest block number
	GetLatestBlockNumber(ctx context.Context) (blockchain.BlockNumber, error)
	// GetBlockByNumber returns the block by number with transactions
	GetBlockByNumber(ctx context.Context, blockNumber blockchain.BlockNumber) (*blockchain.Block, error)
}

var _ Provider = (*provider)(nil)

type provider struct{}

func NewProvider() Provider {
	return &provider{}
}

func (p *provider) GetLatestBlockNumber(ctx context.Context) (blockchain.BlockNumber, error) {
	req, err := getLatestBlockNumberHTTPRequest(ctx)
	if err != nil {
		log.Println("Error creating request:", err)
		return blockchain.InvalidBlockNumber, err
	}
	httpClient := getHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return blockchain.InvalidBlockNumber, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return blockchain.InvalidBlockNumber, err
	}

	var rpcResponse jsonrpc.Response
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return blockchain.InvalidBlockNumber, err
	}
	if rpcResponse.Error != nil {
		log.Println("Error response:", rpcResponse.Error)
		return blockchain.InvalidBlockNumber, err
	}

	var resultHexStr string
	err = json.Unmarshal(rpcResponse.Result, &resultHexStr)
	if err != nil {
		log.Println("Error unmarshalling result:", err)
		return blockchain.InvalidBlockNumber, err
	}
	blockNumber := blockchain.NewBlockNumberBuilder().FromHexString(resultHexStr).Value()
	return blockNumber, nil
}

func getLatestBlockNumberHTTPRequest(ctx context.Context) (*http.Request, error) {
	request := jsonrpc.NewRequest("eth_blockNumber", json.RawMessage("[]"))
	payload, err := json.Marshal(request)
	if err != nil {
		log.Println("Error marshaling request:", err)
		return nil, err
	}

	return http.NewRequestWithContext(ctx, http.MethodPost, cloudFlareRpcURL, bytes.NewBuffer(payload))
}

func (p *provider) GetBlockByNumber(ctx context.Context, blockNumber blockchain.BlockNumber) (*blockchain.Block, error) {
	req, err := getBlockByNumberHTTPRequest(ctx, blockNumber)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}
	httpClient := getHttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResponse jsonrpc.Response
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return nil, err
	}
	if rpcResponse.Error != nil {
		log.Println("Error response:", rpcResponse.Error)
		return nil, err
	}

	var block blockchain.Block
	if err := json.Unmarshal(rpcResponse.Result, &block); err != nil {
		log.Println("Error unmarshalling block:", err)
		return nil, err
	}
	return &block, nil
}

func getBlockByNumberHTTPRequest(ctx context.Context, blockNumber blockchain.BlockNumber) (*http.Request, error) {
	paramsStr := fmt.Sprintf(`["0x%x", true]`, blockNumber)
	request := jsonrpc.NewRequest("eth_getBlockByNumber", json.RawMessage(paramsStr))
	payload, err := json.Marshal(request)
	if err != nil {
		log.Println("Error marshaling request:", err)
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodPost, cloudFlareRpcURL, bytes.NewBuffer(payload))
}

func getHttpClient() *http.Client {
	return &http.Client{
		Timeout: httpTimeout,
	}
}
