package main

import (
	"encoding/hex"

	"github.com/umbracle/fastrlp"
)

// sealHash returns the hash of a block.
func sealHash(block map[string]any) ([]byte, error) {
	parentHash, err := hex.DecodeString(block["parentHash"].(string))
	if err != nil {
		return nil, err
	}
	sha3Uncles, err := hex.DecodeString(block["sha3Uncles"].(string))
	if err != nil {
		return nil, err
	}
	miner, err := hex.DecodeString(block["miner"].(string))
	if err != nil {
		return nil, err
	}
	stateRoot, err := hex.DecodeString(block["stateRoot"].(string))
	if err != nil {
		return nil, err
	}
	transactionsRoot, err := hex.DecodeString(block["transactionsRoot"].(string))
	if err != nil {
		return nil, err
	}
	receiptsRoot, err := hex.DecodeString(block["receiptsRoot"].(string))
	if err != nil {
		return nil, err
	}
	logsBloom, err := hex.DecodeString(block["logsBloom"].(string))
	if err != nil {
		return nil, err
	}
	difficulty, err := hex.DecodeString(block["difficulty"].(string))
	if err != nil {
		return nil, err
	}
	number, err := hex.DecodeString(block["number"].(string))
	if err != nil {
		return nil, err
	}
	gasLimit, err := hex.DecodeString(block["gasLimit"].(string))
	if err != nil {
		return nil, err
	}
	gasUsed, err := hex.DecodeString(block["gasUsed"].(string))
	if err != nil {
		return nil, err
	}
	timestamp, err := hex.DecodeString(block["timestamp"].(string))
	if err != nil {
		return nil, err
	}
	extraData, err := hex.DecodeString(block["extraData"].(string))
	if err != nil {
		return nil, err
	}

	var baseFeePerGas []byte
	if val, ok := block["baseFeePerGas"]; ok {
		baseFeePerGas, err = hex.DecodeString(val.(string))
	}
	if err != nil {
		return nil, err
	}

	var withdrawalsRoot []byte
	if val, ok := block["withdrawalsRoot"]; ok {
		withdrawalsRoot, err = hex.DecodeString(val.(string))
	}
	if err != nil {
		return nil, err
	}

	rlp := &fastrlp.Arena{}
	enc := rlp.NewArray()
	enc.Set(rlp.NewBytes(parentHash))
	enc.Set(rlp.NewBytes(sha3Uncles))
	enc.Set(rlp.NewBytes(miner))
	enc.Set(rlp.NewBytes(stateRoot))
	enc.Set(rlp.NewBytes(transactionsRoot))
	enc.Set(rlp.NewBytes(receiptsRoot))
	enc.Set(rlp.NewBytes(logsBloom))
	enc.Set(rlp.NewBytes(difficulty))
	enc.Set(rlp.NewBytes(number))
	enc.Set(rlp.NewBytes(gasLimit))
	enc.Set(rlp.NewBytes(gasUsed))
	enc.Set(rlp.NewBytes(timestamp))
	enc.Set(rlp.NewBytes(extraData))

	if baseFeePerGas != nil {
		enc.Set(rlp.NewBytes(baseFeePerGas))
	}
	if withdrawalsRoot != nil {
		enc.Set(rlp.NewBytes(withdrawalsRoot))
	}

	encoded := enc.MarshalTo(nil)
	keccak := fastrlp.NewKeccak256()
	return keccak.Sum(encoded), nil
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
