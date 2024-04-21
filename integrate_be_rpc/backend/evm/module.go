package evm

import (
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

func (m *EvmBackend) GetEvmModuleParams() (*evmtypes.Params, error) {
	res, err := m.queryClient.EvmQueryClient.Params(m.ctx, &evmtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	return &res.Params, nil
}
