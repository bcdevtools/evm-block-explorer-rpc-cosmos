package message_involves_extractors

import (
	"encoding/hex"
	berpc "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc"
	berpctypes "github.com/bcdevtools/block-explorer-rpc-cosmos/be_rpc/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
)

func RegisterMessageInvolvesExtractorsForErc20() {
	berpc.RegisterMessageInvolversExtractor(&erc20types.MsgConvertCoin{}, ExtractFromErc20MsgConvertCoin)
	berpc.RegisterMessageInvolversExtractor(&erc20types.MsgConvertERC20{}, ExtractFromErc20MsgConvertERC20)
}

func ExtractFromErc20MsgConvertCoin(sdkMsg sdk.Msg, _ *tx.Tx, _ tmtypes.Tx, _ client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*erc20types.MsgConvertCoin)

	res = berpctypes.NewMessageInvolversResult()

	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Sender)
	if strings.HasPrefix(msg.Receiver, "0x") {
		if bz, err := hex.DecodeString(msg.Receiver[2:]); err == nil {
			res.AddGenericInvolvers(berpctypes.MessageInvolvers, sdk.AccAddress(bz).String())
		}
	}

	return
}

func ExtractFromErc20MsgConvertERC20(sdkMsg sdk.Msg, _ *tx.Tx, _ tmtypes.Tx, _ client.Context) (res berpctypes.MessageInvolversResult, err error) {
	msg := sdkMsg.(*erc20types.MsgConvertERC20)

	res = berpctypes.NewMessageInvolversResult()

	if strings.HasPrefix(msg.Sender, "0x") {
		if bz, err := hex.DecodeString(msg.Sender[2:]); err == nil {
			res.AddGenericInvolvers(berpctypes.MessageInvolvers, sdk.AccAddress(bz).String())
		}
	}
	res.AddGenericInvolvers(berpctypes.MessageInvolvers, msg.Receiver)

	return
}
