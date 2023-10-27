package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/multiformats/go-multiaddr"
	badgerds "github.com/sourcenetwork/defradb/datastore/badger/v4"
	"github.com/sourcenetwork/defradb/db"
	"github.com/sourcenetwork/defradb/net"

	"github.com/sourcenetwork/defradb-blockchain-indexer/evm"
)

var (
	rpcFlag = flag.String("rpc", "", "rpc address")
	p2pFlag = flag.String("p2p", "/ip4/0.0.0.0/tcp/9181", "p2p bind address")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := badgerds.Options{
		Options: badger.DefaultOptions("").WithInMemory(true),
	}
	rootstore, err := badgerds.NewDatastore("", &opts)
	if err != nil {
		log.Fatalf("failed to create datastore: %v", err)
	}
	database, err := db.NewDB(ctx, rootstore)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer database.Close()

	addr, err := multiaddr.NewMultiaddr(*p2pFlag)
	if err != nil {
		log.Fatalf("failed to parse multiaddr: %v", err)
	}
	node, err := net.NewNode(ctx, database, net.WithListenAddrs(addr))
	if err != nil {
		log.Fatalf("failed to create node: %v", err)
	}
	indexer, err := evm.NewIndexer(ctx, node, *rpcFlag)
	if err != nil {
		log.Fatalf("failed to create indexer: %v", err)
	}
	go func() {
		if err := indexer.Start(ctx, 5); err != nil {
			log.Fatalf("indexer shutdown: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
}
