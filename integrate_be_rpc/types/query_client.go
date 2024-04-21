package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// QueryClient defines a gRPC Client
type QueryClient struct {
	tx.ServiceClient
	BankQueryClient banktypes.QueryClient
	EvmQueryClient  evmtypes.QueryClient
}

// NewQueryClient creates a new gRPC query clients
func NewQueryClient(clientCtx client.Context) *QueryClient {
	queryClient := &QueryClient{
		ServiceClient:   tx.NewServiceClient(clientCtx),
		BankQueryClient: banktypes.NewQueryClient(clientCtx),
		EvmQueryClient:  evmtypes.NewQueryClient(clientCtx),
	}
	return queryClient
}
