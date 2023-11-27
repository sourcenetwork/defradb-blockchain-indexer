package main

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/umbracle/fastrlp"
)

var rlp = &fastrlp.Arena{}

const (
	LegacyTxType     = 0x00
	AccessListTxType = 0x01
	DynamicFeeTxType = 0x02
	BlobTxType       = 0x03
)

func sealLegacyTx(transaction map[string]any) ([]byte, error) {
	nonce, err := hex.DecodeString(transaction["nonce"].(string))
	if err != nil {
		return nil, err
	}
	gasPrice, err := hex.DecodeString(transaction["gasPrice"].(string))
	if err != nil {
		return nil, err
	}
	gasUsed, err := hex.DecodeString(transaction["gasUsed"].(string))
	if err != nil {
		return nil, err
	}
	to, err := hex.DecodeString(transaction["to"].(string))
	if err != nil {
		return nil, err
	}
	value, err := hex.DecodeString(transaction["value"].(string))
	if err != nil {
		return nil, err
	}

	var input []byte
	if val, ok := transaction["input"].(string); ok {
		input, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var chainId []byte
	if val, ok := transaction["chainId"].(string); ok {
		chainId, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	enc := rlp.NewArray()
	enc.Set(rlp.NewBytes(nonce))
	enc.Set(rlp.NewBytes(gasPrice))
	enc.Set(rlp.NewBytes(gasUsed))
	enc.Set(rlp.NewBytes(to))
	enc.Set(rlp.NewBytes(value))
	enc.Set(rlp.NewBytes(input))

	if chainId != nil {
		enc.Set(rlp.NewBytes(chainId))
		enc.Set(rlp.NewUint(0))
		enc.Set(rlp.NewUint(0))
	}

	keccak := fastrlp.NewKeccak256()
	if _, err := keccak.Write(enc.MarshalTo(nil)); err != nil {
		return nil, err
	}
	return keccak.Sum(nil), nil
}

func sealAccessListTx(transaction map[string]any) ([]byte, error) {
	nonce, err := hex.DecodeString(transaction["nonce"].(string))
	if err != nil {
		return nil, err
	}
	gasPrice, err := hex.DecodeString(transaction["gasPrice"].(string))
	if err != nil {
		return nil, err
	}
	gasUsed, err := hex.DecodeString(transaction["gasUsed"].(string))
	if err != nil {
		return nil, err
	}
	to, err := hex.DecodeString(transaction["to"].(string))
	if err != nil {
		return nil, err
	}
	value, err := hex.DecodeString(transaction["value"].(string))
	if err != nil {
		return nil, err
	}

	var input []byte
	if val, ok := transaction["input"].(string); ok {
		input, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var chainId []byte
	if val, ok := transaction["chainId"].(string); ok {
		chainId, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	enc := rlp.NewArray()
	enc.Set(rlp.NewBytes(chainId))
	enc.Set(rlp.NewBytes(nonce))
	enc.Set(rlp.NewBytes(gasPrice))
	enc.Set(rlp.NewBytes(gasUsed))
	enc.Set(rlp.NewBytes(to))
	enc.Set(rlp.NewBytes(value))
	enc.Set(rlp.NewBytes(input))
	// TODO encode access list

	keccak := fastrlp.NewKeccak256()
	if _, err := keccak.Write([]byte{AccessListTxType}); err != nil {
		return nil, err
	}
	if _, err := keccak.Write(enc.MarshalTo(nil)); err != nil {
		return nil, err
	}
	return keccak.Sum(nil), nil
}

func sealDynamicFeeTx(transaction map[string]any) ([]byte, error) {
	nonce, err := hex.DecodeString(transaction["nonce"].(string))
	if err != nil {
		return nil, err
	}
	gasPrice, err := hex.DecodeString(transaction["gasPrice"].(string))
	if err != nil {
		return nil, err
	}
	gasUsed, err := hex.DecodeString(transaction["gasUsed"].(string))
	if err != nil {
		return nil, err
	}
	to, err := hex.DecodeString(transaction["to"].(string))
	if err != nil {
		return nil, err
	}
	value, err := hex.DecodeString(transaction["value"].(string))
	if err != nil {
		return nil, err
	}

	var input []byte
	if val, ok := transaction["input"].(string); ok {
		input, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var chainId []byte
	if val, ok := transaction["chainId"].(string); ok {
		chainId, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var maxPriorityFeePerGas []byte
	if val, ok := transaction["maxPriorityFeePerGas"].(string); ok {
		maxPriorityFeePerGas, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	enc := rlp.NewArray()
	enc.Set(rlp.NewBytes(chainId))
	enc.Set(rlp.NewBytes(nonce))
	enc.Set(rlp.NewBytes(maxPriorityFeePerGas))
	enc.Set(rlp.NewBytes(gasPrice))
	enc.Set(rlp.NewBytes(gasUsed))
	enc.Set(rlp.NewBytes(to))
	enc.Set(rlp.NewBytes(value))
	enc.Set(rlp.NewBytes(input))
	// TODO encode access list

	keccak := fastrlp.NewKeccak256()
	if _, err := keccak.Write([]byte{DynamicFeeTxType}); err != nil {
		return nil, err
	}
	if _, err := keccak.Write(enc.MarshalTo(nil)); err != nil {
		return nil, err
	}
	return keccak.Sum(nil), nil
}

func sealBlobTx(transaction map[string]any) ([]byte, error) {
	nonce, err := hex.DecodeString(transaction["nonce"].(string))
	if err != nil {
		return nil, err
	}
	gasPrice, err := hex.DecodeString(transaction["gasPrice"].(string))
	if err != nil {
		return nil, err
	}
	gasUsed, err := hex.DecodeString(transaction["gasUsed"].(string))
	if err != nil {
		return nil, err
	}
	to, err := hex.DecodeString(transaction["to"].(string))
	if err != nil {
		return nil, err
	}
	value, err := hex.DecodeString(transaction["value"].(string))
	if err != nil {
		return nil, err
	}

	var input []byte
	if val, ok := transaction["input"].(string); ok {
		input, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var chainId []byte
	if val, ok := transaction["chainId"].(string); ok {
		chainId, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var maxPriorityFeePerGas []byte
	if val, ok := transaction["maxPriorityFeePerGas"].(string); ok {
		maxPriorityFeePerGas, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	var maxFeePerBlobGas []byte
	if val, ok := transaction["maxFeePerBlobGas"].(string); ok {
		maxFeePerBlobGas, err = hex.DecodeString(val)
	}
	if err != nil {
		return nil, err
	}

	blobHashes := rlp.NewArray()
	for _, blob := range transaction["blobVersionedHashes"].([]string) {
		hash, err := hex.DecodeString(blob)
		if err != nil {
			return nil, err
		}
		blobHashes.Set(rlp.NewBytes(hash))
	}

	enc := rlp.NewArray()
	enc.Set(rlp.NewBytes(chainId))
	enc.Set(rlp.NewBytes(nonce))
	enc.Set(rlp.NewBytes(maxPriorityFeePerGas))
	enc.Set(rlp.NewBytes(gasPrice))
	enc.Set(rlp.NewBytes(gasUsed))
	enc.Set(rlp.NewBytes(to))
	enc.Set(rlp.NewBytes(value))
	enc.Set(rlp.NewBytes(input))
	// TODO encode access list
	enc.Set(rlp.NewBytes(maxFeePerBlobGas))
	enc.Set(blobHashes)

	keccak := fastrlp.NewKeccak256()
	if _, err := keccak.Write([]byte{BlobTxType}); err != nil {
		return nil, err
	}
	if _, err := keccak.Write(enc.MarshalTo(nil)); err != nil {
		return nil, err
	}
	return keccak.Sum(nil), nil
}

func seal(transaction map[string]any) ([]byte, error) {
	var (
		prefix uint64
		err    error
	)
	if val, ok := transaction["prefix"].(string); ok {
		prefix, err = strconv.ParseUint(val, 0, 8)
	}
	if err != nil {
		return nil, err
	}

	switch prefix {
	case LegacyTxType:
		return sealLegacyTx(transaction)

	case AccessListTxType:
		return sealAccessListTx(transaction)

	case DynamicFeeTxType:
		return sealDynamicFeeTx(transaction)

	case BlobTxType:
		return sealBlobTx(transaction)

	default:
		return nil, fmt.Errorf("invalid transaction type: %d", prefix)
	}
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
