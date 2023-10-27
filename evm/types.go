package evm

type Block struct {
	Number                string        `json:"number"`
	Hash                  string        `json:"hash"`
	ParentHash            string        `json:"parentHash"`
	Nonce                 string        `json:"nonce"`
	Sha3Uncles            string        `json:"sha3Uncles"`
	LogsBloom             string        `json:"logsBloom"`
	TransactionsRoot      string        `json:"transactionsRoot"`
	StateRoot             string        `json:"stateRoot"`
	ReceiptsRoot          string        `json:"receiptsRoot"`
	Miner                 string        `json:"miner"`
	Difficulty            string        `json:"difficulty"`
	TotalDifficulty       string        `json:"totalDifficulty"`
	ExtraData             string        `json:"extraData"`
	Size                  string        `json:"size"`
	GasLimit              string        `json:"gasLimit"`
	GasUsed               string        `json:"gasUsed"`
	Timestamp             string        `json:"timestamp"`
	MixHash               string        `json:"mixHash"`
	BaseFeePerGas         string        `json:"baseFeePerGas"`
	WithdrawalsRoot       string        `json:"withdrawalsRoot"`
	BlobGasUsed           string        `json:"blobGasUsed"`
	ExcessBlobGas         string        `json:"excessBlobGas"`
	ParentBeaconBlockRoot string        `json:"parentBeaconBlockRoot"`
	Transactions          []Transaction `json:"transactions"`
}

func (b *Block) ToMap() map[string]any {
	return map[string]any{
		"number":                b.Number,
		"hash":                  b.Hash,
		"parentHash":            b.ParentHash,
		"nonce":                 b.Nonce,
		"sha3Uncles":            b.Sha3Uncles,
		"logsBloom":             b.LogsBloom,
		"transactionsRoot":      b.TransactionsRoot,
		"stateRoot":             b.StateRoot,
		"receiptsRoot":          b.ReceiptsRoot,
		"miner":                 b.Miner,
		"difficulty":            b.Difficulty,
		"totalDifficulty":       b.TotalDifficulty,
		"extraData":             b.ExtraData,
		"size":                  b.Size,
		"gasLimit":              b.GasLimit,
		"gasUsed":               b.GasUsed,
		"timestamp":             b.Timestamp,
		"mixHash":               b.MixHash,
		"baseFeePerGas":         b.BaseFeePerGas,
		"withdrawalsRoot":       b.WithdrawalsRoot,
		"blobGasUsed":           b.BlobGasUsed,
		"excessBlobGas":         b.ExcessBlobGas,
		"parentBeaconBlockRoot": b.ParentBeaconBlockRoot,
	}
}

type Transaction struct {
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func (t *Transaction) ToMap() map[string]any {
	return map[string]any{
		"from":             t.From,
		"gas":              t.Gas,
		"gasPrice":         t.GasPrice,
		"hash":             t.Hash,
		"input":            t.Input,
		"nonce":            t.Nonce,
		"to":               t.To,
		"transactionIndex": t.TransactionIndex,
		"value":            t.Value,
		"v":                t.V,
		"r":                t.R,
		"s":                t.S,
	}
}
