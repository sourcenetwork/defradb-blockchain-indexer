# Ethereum Transaction Validator

This lens creates a view of validated transactions.

Transactions are sealed using the Keccak256 hash of the RLP encoded fields.

## Building

```bash
tinygo build -o main.wasm -target=wasi main.go
```
