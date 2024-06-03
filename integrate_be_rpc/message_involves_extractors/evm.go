package message_involves_extractors

import (
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/bcdevtools/evm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func RegisterMessageInvolvesExtractorsForEvm(evmBeRpcBackend evm.EvmBackendI) {
	berpc.RegisterMessageInvolversExtractor(&evmtypes.MsgEthereumTx{}, func(msg sdk.Msg, _ *tx.Tx, _ tmtypes.Tx, _ client.Context) (berpctypes.MessageInvolversResult, error) {
		return evmBeRpcBackend.GetEvmTransactionInvolversByHash(
			msg.(*evmtypes.MsgEthereumTx).AsTransaction().Hash(),
		)
	})
}
