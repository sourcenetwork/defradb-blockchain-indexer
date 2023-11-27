package main

import (
	"encoding/hex"

	"github.com/umbracle/fastrlp"
)

var rlp = &fastrlp.Arena{}

// seal returns the hash of a header.
func seal(header map[string]any) ([]byte, error) {
	parentHash, err := hex.DecodeString(header["parentHash"].(string))
	if err != nil {
		return nil, err
	}
	sha3Uncles, err := hex.DecodeString(header["sha3Uncles"].(string))
	if err != nil {
		return nil, err
	}
	miner, err := hex.DecodeString(header["miner"].(string))
	if err != nil {
		return nil, err
	}
	stateRoot, err := hex.DecodeString(header["stateRoot"].(string))
	if err != nil {
		return nil, err
	}
	transactionsRoot, err := hex.DecodeString(header["transactionsRoot"].(string))
	if err != nil {
		return nil, err
	}
	receiptsRoot, err := hex.DecodeString(header["receiptsRoot"].(string))
	if err != nil {
		return nil, err
	}
	logsBloom, err := hex.DecodeString(header["logsBloom"].(string))
	if err != nil {
		return nil, err
	}
	difficulty, err := hex.DecodeString(header["difficulty"].(string))
	if err != nil {
		return nil, err
	}
	number, err := hex.DecodeString(header["number"].(string))
	if err != nil {
		return nil, err
	}
	gasLimit, err := hex.DecodeString(header["gasLimit"].(string))
	if err != nil {
		return nil, err
	}
	gasUsed, err := hex.DecodeString(header["gasUsed"].(string))
	if err != nil {
		return nil, err
	}
	timestamp, err := hex.DecodeString(header["timestamp"].(string))
	if err != nil {
		return nil, err
	}
	extraData, err := hex.DecodeString(header["extraData"].(string))
	if err != nil {
		return nil, err
	}

	var baseFeePerGas []byte
	if val, ok := header["baseFeePerGas"]; ok {
		baseFeePerGas, err = hex.DecodeString(val.(string))
	}
	if err != nil {
		return nil, err
	}

	var withdrawalsRoot []byte
	if val, ok := header["withdrawalsRoot"]; ok {
		withdrawalsRoot, err = hex.DecodeString(val.(string))
	}
	if err != nil {
		return nil, err
	}

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
	if _, err := keccak.Write(enc.MarshalTo(nil)); err != nil {
		return nil, err
	}
	return keccak.Sum(nil), nil
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
