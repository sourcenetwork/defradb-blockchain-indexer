package main

import (
	"bytes"
	"encoding/hex"

	"github.com/umbracle/fastrlp"
)

// validateSeal ensures that the block header matches the hash of the encoded fields
func validateSeal(block map[string]any) (bool, error) {
	parentHash, err := hex.DecodeString(block["parentHash"].(string))
	if err != nil {
		return false, err
	}
	sha3Uncles, err := hex.DecodeString(block["sha3Uncles"].(string))
	if err != nil {
		return false, err
	}
	miner, err := hex.DecodeString(block["miner"].(string))
	if err != nil {
		return false, err
	}
	stateRoot, err := hex.DecodeString(block["stateRoot"].(string))
	if err != nil {
		return false, err
	}
	transactionsRoot, err := hex.DecodeString(block["transactionsRoot"].(string))
	if err != nil {
		return false, err
	}
	receiptsRoot, err := hex.DecodeString(block["receiptsRoot"].(string))
	if err != nil {
		return false, err
	}
	logsBloom, err := hex.DecodeString(block["logsBloom"].(string))
	if err != nil {
		return false, err
	}
	difficulty, err := hex.DecodeString(block["difficulty"].(string))
	if err != nil {
		return false, err
	}
	number, err := hex.DecodeString(block["number"].(string))
	if err != nil {
		return false, err
	}
	gasLimit, err := hex.DecodeString(block["gasLimit"].(string))
	if err != nil {
		return false, err
	}
	gasUsed, err := hex.DecodeString(block["gasUsed"].(string))
	if err != nil {
		return false, err
	}
	timestamp, err := hex.DecodeString(block["timestamp"].(string))
	if err != nil {
		return false, err
	}
	extraData, err := hex.DecodeString(block["extraData"].(string))
	if err != nil {
		return false, err
	}

	var baseFeePerGas []byte
	if val, ok := block["baseFeePerGas"]; ok {
		baseFeePerGas, err = hex.DecodeString(val.(string))
	}
	if err != nil {
		return false, err
	}

	var withdrawalsRoot []byte
	if val, ok := block["withdrawalsRoot"]; ok {
		withdrawalsRoot, err = hex.DecodeString(val.(string))
	}
	if err != nil {
		return false, err
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

	keccak := fastrlp.NewKeccak256()
	seal := keccak.Sum(enc.MarshalTo(nil))
	hash, err := hex.DecodeString(block["hash"].(string))
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash, seal), nil
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
