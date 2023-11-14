# Ethereum Lens

This lens creates a view of sealed blocks.

## Transform

Blocks are sealed using the Keccak256 hash of the RLP encoded fields.

## Building

```bash
tinygo build -o main.wasm -target=wasi main.go
```
