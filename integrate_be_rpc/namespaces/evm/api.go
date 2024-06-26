package evm

import (
	"fmt"
	ieberpcbackend "github.com/bcdevtools/evm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/tendermint/tendermint/libs/log"
)

// RPC namespaces and API version
const (
	DymEvmBlockExplorerNamespace = "evm"

	ApiVersion = "1.0"
)

// API is the EVM Block Explorer JSON-RPC.
type API struct {
	ctx     *server.Context
	logger  log.Logger
	backend ieberpcbackend.EvmBackendI
}

// NewEvmBeAPI creates an instance of the EVM Block Explorer API.
func NewEvmBeAPI(
	ctx *server.Context,
	backend ieberpcbackend.EvmBackendI,
) *API {
	return &API{
		ctx:     ctx,
		logger:  ctx.Logger.With("api", "evm"),
		backend: backend,
	}
}

func (api *API) Echo(text string) string {
	api.logger.Debug("evm_echo")
	return fmt.Sprintf("hello \"%s\" from EVM Block Explorer API", text)
}
