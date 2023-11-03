package ethereum

import (
	"context"
	_ "embed"
	"log"

	"github.com/sourcenetwork/defradb/client"
)

//go:embed schema.gql
var schema string

type Indexer struct {
	database     client.DB
	rpc          *RpcClient
	chainID      string
	blocks       client.Collection
	transactions client.Collection
	withdrawals  client.Collection
}

func NewIndexer(ctx context.Context, database client.DB, url string) (*Indexer, error) {
	_, err := database.AddSchema(ctx, schema)
	if err != nil {
		return nil, err
	}
	blocks, err := database.GetCollectionByName(ctx, "EthereumBlock")
	if err != nil {
		return nil, err
	}
	transactions, err := database.GetCollectionByName(ctx, "EthereumTransaction")
	if err != nil {
		return nil, err
	}
	withdrawals, err := database.GetCollectionByName(ctx, "EthereumWithdrawal")
	if err != nil {
		return nil, err
	}
	rpc, err := NewRpcClient(ctx, url)
	if err != nil {
		return nil, err
	}
	chainID, err := rpc.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return &Indexer{
		database:     database,
		rpc:          rpc,
		chainID:      chainID,
		blocks:       blocks,
		transactions: transactions,
		withdrawals:  withdrawals,
	}, nil
}

// Start runs the indexer with the given number of processes.
func (idx *Indexer) Start(ctx context.Context, count int) error {
	blockCh := make(chan uint64)
	defer close(blockCh)

	log.Printf("Starting indexer with %d processes...", count)
	for i := 0; i < count; i++ {
		go func() { idx.indexBlockProcess(ctx, blockCh) }()
	}
	return idx.indexHeadProcess(ctx, blockCh, nil)
}

func (idx *Indexer) indexBlock(ctx context.Context, block map[string]any) error {
	txn, err := idx.database.NewTxn(ctx, false)
	defer txn.Discard(ctx)

	var transactions []map[string]any
	if val, ok := block["transactions"].([]map[string]any); ok {
		transactions = val
	}
	var withdrawals []map[string]any
	if val, ok := block["withdrawals"].([]map[string]any); ok {
		withdrawals = val
	}

	delete(block, "transactions")
	delete(block, "withdrawals")

	block["chainID"] = idx.chainID
	blockDoc, err := client.NewDocFromMap(block)
	if err != nil {
		return err
	}
	err = idx.blocks.WithTxn(txn).Create(ctx, blockDoc)
	if err != nil {
		return err
	}

	for _, withdrawal := range withdrawals {
		withdrawal["chainID"] = idx.chainID
		withdrawal["blockHash"] = block["hash"]
		withdrawal["blockNumber"] = block["number"]
		withdrawalDoc, err := client.NewDocFromMap(withdrawal)
		if err != nil {
			return err
		}
		err = idx.transactions.WithTxn(txn).Create(ctx, withdrawalDoc)
		if err != nil {
			return err
		}
	}

	for _, transaction := range transactions {
		transaction["chainID"] = idx.chainID
		transactionDoc, err := client.NewDocFromMap(transaction)
		if err != nil {
			return err
		}
		err = idx.transactions.WithTxn(txn).Create(ctx, transactionDoc)
		if err != nil {
			return err
		}
	}
	return txn.Commit(ctx)
}
