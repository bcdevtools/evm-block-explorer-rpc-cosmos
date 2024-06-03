package message_parsers

import (
	"fmt"
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	berpcutils "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/utils"
	"github.com/bcdevtools/evm-block-explorer-rpc-cosmos/integrate_be_rpc/backend/evm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
	"strings"
)

func RegisterMessageParsersForEvm(evmBeRpcBackend evm.EvmBackendI) {
	berpc.RegisterMessageParser(&erc20types.MsgConvertCoin{}, func(sdkMsg sdk.Msg, _ uint, _ *tx.Tx, _ *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
		msg := sdkMsg.(*erc20types.MsgConvertCoin)

		res = berpctypes.GenericBackendResponse{
			"transfer": map[string]any{
				"from": []string{msg.Sender},
				"to": []map[string]any{
					{
						"address": strings.ToLower(msg.Receiver),
						"amount":  berpcutils.CoinsToMap(msg.Coin),
					},
				},
			},
		}

		berpctypes.NewFriendlyResponseContentBuilder().
			WriteAddress(msg.Sender).
			WriteText(" converts ").
			WriteCoins(sdk.Coins{msg.Coin}, evmBeRpcBackend.GetBankDenomsMetadata(sdk.Coins{msg.Coin})).
			WriteText(" to ERC-20 tokens and send to ").
			WriteAddress(strings.ToLower(msg.Receiver)).
			BuildIntoResponse(res)

		return
	})

	berpc.RegisterMessageParser(&erc20types.MsgConvertERC20{}, func(sdkMsg sdk.Msg, _ uint, _ *tx.Tx, _ *sdk.TxResponse) (res berpctypes.GenericBackendResponse, err error) {
		msg := sdkMsg.(*erc20types.MsgConvertERC20)

		res = berpctypes.GenericBackendResponse{
			"transfer": map[string]any{
				"from": []string{strings.ToLower(msg.Sender)},
				"to": []map[string]any{
					{
						"address": msg.Receiver,
						"amount":  berpcutils.CoinsToMap(sdk.NewCoin(fmt.Sprintf("erc20/%s", msg.ContractAddress), msg.Amount)),
					},
				},
			},
		}

		rb := berpctypes.NewFriendlyResponseContentBuilder().
			WriteAddress(strings.ToLower(msg.Sender)).
			WriteText(" converts ERC-20 token (contract ")

		if contractAddr := common.HexToAddress(msg.ContractAddress); !berpcutils.IsZeroEvmAddress(contractAddr) {
			if contractInfo, err := evmBeRpcBackend.GetErc20ContractInfo(contractAddr); err == nil {
				if contractName, found := contractInfo["name"]; found {
					if contractNameStr, ok := contractName.(string); ok && contractNameStr != "" {
						rb.WriteText(contractNameStr).WriteText(" ")
					}
				}
			}
		}

		rb.WriteAddress(strings.ToLower(msg.ContractAddress)).WriteText(") back to coins and send to ").
			WriteAddress(msg.Receiver).
			BuildIntoResponse(res)

		return
	})
}
