## Integrate EVM Block Explorer Json-RPC module into your chain

The following methods must be called:
```go
config.EnsureRoot(home, config.DefaultBeJsonRpcConfig())
// in root.go
```
```go
config.AddBeJsonRpcFlags(rootCmd)
// in start.go
```
```go
server.StartEvmBeJsonRPC(...)
// in start.go
```