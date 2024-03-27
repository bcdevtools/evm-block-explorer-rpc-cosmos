package integrate_be_rpc

import (
	"context"
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpcbackend "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/backend"
	berpccfg "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/config"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcserver "github.com/bcdevtools/block-explorer-rpc-cosmos/server"
	evmberpcbackend "github.com/bcdevtools/evm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	evmbeapi "github.com/bcdevtools/evm-block-explorer-rpc-cosmos/integrate_be_rpc/namespaces/evm"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	evmostypes "github.com/evmos/evmos/v12/types"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	rpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

type FuncRegister func(evmberpcbackend.EvmBackendI)

func StartEvmBeJsonRPC(
	ctx *server.Context,
	clientCtx client.Context,
	chainId string,
	beRpcCfg berpccfg.BeJsonRpcConfig,
	evmTxIndexer evmostypes.EVMTxIndexer,
	externalServicesModifierFunc func(berpctypes.ExternalServices) berpctypes.ExternalServices,
	apiNamespaceRegisterFunc FuncRegister,
	customInterceptorCreationFunc func(berpcbackend.BackendI, evmberpcbackend.EvmBackendI) berpcbackend.RequestInterceptor,
	tmRPCAddr, tmEndpoint string,
) (serverCloseDeferFunc func(), err error) {
	if err := beRpcCfg.Validate(); err != nil {
		return nil, err
	}

	clientCtx = clientCtx.WithChainID(chainId)

	externalServices := berpctypes.ExternalServices{
		ChainType:    berpctypes.ChainTypeEvm,
		EvmTxIndexer: beJsonRpcEvmIndexerCompatible{indexer: evmTxIndexer},
	}
	if externalServicesModifierFunc != nil {
		externalServices = externalServicesModifierFunc(externalServices)
	}

	evmBeRpcBackend := evmberpcbackend.NewEvmBackend(ctx, ctx.Logger, clientCtx, externalServices)

	berpc.RegisterAPINamespace(evmbeapi.DymEvmBlockExplorerNamespace, func(ctx *server.Context,
		_ client.Context,
		_ *rpcclient.WSClient,
		_ map[string]berpctypes.MessageParser,
		_ map[string]berpctypes.MessageInvolversExtractor,
		_ func(berpcbackend.BackendI) berpcbackend.RequestInterceptor,
		_ berpctypes.ExternalServices,
	) []rpc.API {
		return []rpc.API{
			{
				Namespace: evmbeapi.DymEvmBlockExplorerNamespace,
				Version:   evmbeapi.ApiVersion,
				Service:   evmbeapi.NewEvmBeAPI(ctx, evmBeRpcBackend),
				Public:    true,
			},
		}
	}, false)

	if apiNamespaceRegisterFunc != nil {
		apiNamespaceRegisterFunc(evmBeRpcBackend)
	}

	// register message parsers & message involvers extractor

	berpc.RegisterMessageInvolversExtractor(&evmtypes.MsgEthereumTx{}, func(msg sdk.Msg, _ *tx.Tx, _ tmtypes.Tx, _ client.Context) (berpctypes.MessageInvolversResult, error) {
		return evmBeRpcBackend.GetEvmTransactionInvolversByHash(
			msg.(*evmtypes.MsgEthereumTx).AsTransaction().Hash(),
		)
	})

	var interceptorCreationFunc func(berpcbackend.BackendI) berpcbackend.RequestInterceptor
	if customInterceptorCreationFunc != nil {
		interceptorCreationFunc = func(backend berpcbackend.BackendI) berpcbackend.RequestInterceptor {
			return customInterceptorCreationFunc(backend, evmBeRpcBackend)
		}
	}

	beJsonRpcHttpSrv, beJsonRpcHttpSrvDone, err := berpcserver.StartBeJsonRPC(
		ctx, clientCtx, tmRPCAddr, tmEndpoint,
		beRpcCfg,
		interceptorCreationFunc,
		externalServices,
	)
	if err != nil {
		return nil, err
	}

	return func() {
		shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn()
		if err := beJsonRpcHttpSrv.Shutdown(shutdownCtx); err != nil {
			ctx.Logger.Error("EVM Block Explorer Json-RPC HTTP server shutdown produced a warning", "error", err.Error())
		} else {
			ctx.Logger.Info("EVM Block Explorer Json-RPC HTTP server shut down, waiting 5 sec")
			select {
			case <-time.Tick(5 * time.Second):
			case <-beJsonRpcHttpSrvDone:
			}
		}
	}, nil
}

var _ berpctypes.TxResultForExternal = beJsonRpcEvmIndexerTxResultCompatible{}

type beJsonRpcEvmIndexerTxResultCompatible struct {
	txResult evmostypes.TxResult
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetHeight() int64 {
	return b.txResult.Height
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetTxIndex() uint32 {
	return b.txResult.TxIndex
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetMsgIndex() uint32 {
	return b.txResult.MsgIndex
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetEthTxIndex() int32 {
	return b.txResult.EthTxIndex
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetFailed() bool {
	return b.txResult.Failed
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetGasUsed() uint64 {
	return b.txResult.GasUsed
}

func (b beJsonRpcEvmIndexerTxResultCompatible) GetCumulativeGasUsed() uint64 {
	return b.txResult.CumulativeGasUsed
}

var _ berpctypes.ExpectedEVMTxIndexer = beJsonRpcEvmIndexerCompatible{}

type beJsonRpcEvmIndexerCompatible struct {
	indexer evmostypes.EVMTxIndexer
}

func (b beJsonRpcEvmIndexerCompatible) GetIndexer() any {
	return b.indexer
}

func (b beJsonRpcEvmIndexerCompatible) LastIndexedBlock() (int64, error) {
	return b.LastIndexedBlock()
}

func (b beJsonRpcEvmIndexerCompatible) GetByTxHashForExternal(hash common.Hash) (berpctypes.TxResultForExternal, error) {
	txResult, err := b.indexer.GetByTxHash(hash)
	if err != nil {
		return nil, err
	}
	if txResult == nil {
		return nil, nil
	}
	return beJsonRpcEvmIndexerTxResultCompatible{txResult: *txResult}, nil
}
