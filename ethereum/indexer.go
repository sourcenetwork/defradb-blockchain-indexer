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
	database client.DB
	rpc      *RpcClient
	headers  client.Collection
}

func NewIndexer(ctx context.Context, database client.DB, url string) (*Indexer, error) {
	_, err := database.AddSchema(ctx, schema)
	if err != nil {
		return nil, err
	}
	headers, err := database.GetCollectionByName(ctx, "EthereumHeader")
	if err != nil {
		return nil, err
	}
	rpc, err := NewRpcClient(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		database: database,
		rpc:      rpc,
		headers:  headers,
	}, nil
}

// Start runs the indexer with the given number of processes.
func (idx *Indexer) Start(ctx context.Context, count int) error {
	blockCh := make(chan uint64)
	defer close(blockCh)

	log.Printf("Starting indexer with %d processes...", count)
	for i := 0; i < count; i++ {
		go func() { idx.indexHeaderProcess(ctx, blockCh) }()
	}
	return idx.indexChainProcess(ctx, blockCh, nil)
}

func (idx *Indexer) indexHeader(ctx context.Context, header map[string]any) error {
	txn, err := idx.database.NewTxn(ctx, false)
	defer txn.Discard(ctx)

	delete(header, "hash")
	delete(header, "transactions")
	delete(header, "withdrawals")
	delete(header, "uncles")
	delete(header, "mixHash")
	delete(header, "nonce")
	delete(header, "totalDifficulty")
	delete(header, "size")

	headerDoc, err := client.NewDocFromMap(header)
	if err != nil {
		return err
	}
	err = idx.headers.WithTxn(txn).Create(ctx, headerDoc)
	if err != nil {
		return err
	}
	return txn.Commit(ctx)
}
