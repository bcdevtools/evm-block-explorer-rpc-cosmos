package evm

import (
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
)

func (m *EvmBackend) GetEvmModuleParams() (*evmtypes.Params, error) {
	res, err := m.queryClient.EvmQueryClient.Params(m.ctx, &evmtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	return &res.Params, nil
}

func (m *EvmBackend) GetErc20ModuleParams() (*erc20types.Params, error) {
	res, err := m.queryClient.Erc20QueryClient.Params(m.ctx, &erc20types.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	return &res.Params, nil
}
