//go:build !be_json_rpc_evm

package evm

import (
	"github.com/bcdevtools/integrate-block-explorer-rpc-cosmos/integrate_be_rpc/compatible"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/tendermint/tendermint/libs/log"
)

/**
This file is used to get rid of compile error in IDE.
*/

var _ EvmBackendI = (*EvmBackend)(nil)

type EvmBackendI interface {
}

type EvmBackend struct {
}

// NewEvmBackend creates a new EvmBackend instance for EVM Block Explorer.
// This method is for get rid of build error in IDE in final chains.
func NewEvmBackend(
	ctx *server.Context,
	logger log.Logger,
	clientCtx client.Context,
) *EvmBackend {
	compatible.PanicInvalidBuildTag()
	return nil
}
