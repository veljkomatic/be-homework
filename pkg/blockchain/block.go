package blockchain

// Block represents a block in the blockchain
type Block struct {
	Number           string         `json:"number,omitempty"`
	Hash             string         `json:"hash,omitempty"`
	ParentHash       string         `json:"parentHash,omitempty"`
	Nonce            string         `json:"nonce,omitempty"`
	Sha3Uncles       string         `json:"sha3Uncles,omitempty"`
	LogsBloom        string         `json:"logsBloom,omitempty"`
	TransactionsRoot string         `json:"transactionsRoot,omitempty"`
	StateRoot        string         `json:"stateRoot,omitempty"`
	Miner            string         `json:"miner,omitempty"`
	Difficulty       string         `json:"difficulty,omitempty"`
	TotalDifficulty  string         `json:"totalDifficulty,omitempty"`
	ExtraData        string         `json:"extraData,omitempty"`
	Size             string         `json:"size,omitempty"`
	GasLimit         string         `json:"gasLimit,omitempty"`
	GasUsed          string         `json:"gasUsed,omitempty"`
	Timestamp        string         `json:"timestamp,omitempty"`
	Transactions     []*Transaction `json:"transactions,omitempty"`
	Uncles           []string       `json:"uncles,omitempty"`
}

// Transaction represents a transaction in the blockchain
type Transaction struct {
	Hash             string `json:"hash,omitempty"`
	Nonce            string `json:"nonce,omitempty"`
	BlockHash        string `json:"blockHash,omitempty"`
	BlockNumber      string `json:"blockNumber,omitempty"`
	TransactionIndex string `json:"transactionIndex,omitempty"`
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
	Value            string `json:"value,omitempty"`
	GasPrice         string `json:"gasPrice,omitempty"`
	Gas              string `json:"gas,omitempty"`
	Input            string `json:"input,omitempty"`
}
