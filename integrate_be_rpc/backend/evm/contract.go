package evm

import (
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

func (m *EvmBackend) GetContractCode(contractAddress common.Address) ([]byte, error) {
	res, err := m.queryClient.EvmQueryClient.Code(m.ctx, &evmtypes.QueryCodeRequest{Address: contractAddress.String()})
	if err != nil {
		return nil, err
	}
	return res.Code, nil
}
