package evm

import (
	"context"
	_ "embed"
	"errors"
	"log"

	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/defradb/db"
)

//go:embed schema.gql
var schema string

type Indexer struct {
	database    client.DB
	rpc         *RpcClient
	chainDocKey client.DocKey
}

func NewIndexer(ctx context.Context, database client.DB, url string) (*Indexer, error) {
	_, err := database.AddSchema(ctx, schema)
	if err != nil {
		return nil, err
	}
	rpc := NewRpcClient(url)

	chainID, err := rpc.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	chainDoc, err := client.NewDocFromMap(map[string]any{"id": chainID})
	if err != nil {
		return nil, err
	}
	chainDocKey, err := chainDoc.GenerateDocKey()
	if err != nil {
		return nil, err
	}
	chains, err := database.GetCollectionByName(ctx, "Chain")
	if err != nil {
		return nil, err
	}
	err = chains.Create(ctx, chainDoc)
	if err != nil && !errors.Is(err, db.ErrDocumentAlreadyExists) {
		return nil, err
	}

	return &Indexer{
		database:    database,
		rpc:         rpc,
		chainDocKey: chainDocKey,
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
	return idx.indexHeadProcess(ctx, blockCh)
}

func (idx *Indexer) indexBlock(ctx context.Context, block *Block) error {
	txn, err := idx.database.NewTxn(ctx, false)
	defer txn.Discard(ctx)

	blocks, err := idx.database.GetCollectionByName(ctx, "Block")
	if err != nil {
		return err
	}
	transactions, err := idx.database.GetCollectionByName(ctx, "Transaction")
	if err != nil {
		return err
	}

	blockDoc, err := client.NewDocFromMap(block.ToMap())
	if err != nil {
		return err
	}
	err = blockDoc.Set("chain", idx.chainDocKey.String())
	if err != nil {
		return err
	}
	err = blocks.WithTxn(txn).Create(ctx, blockDoc)
	if err != nil {
		return err
	}
	blockKey, err := blockDoc.GenerateDocKey()
	if err != nil {
		return err
	}

	for _, transaction := range block.Transactions {
		transactionDoc, err := client.NewDocFromMap(transaction.ToMap())
		if err != nil {
			return err
		}
		err = transactionDoc.Set("block", blockKey.String())
		if err != nil {
			return err
		}
		err = transactions.WithTxn(txn).Create(ctx, transactionDoc)
		if err != nil {
			return err
		}
	}
	return txn.Commit(ctx)
}
