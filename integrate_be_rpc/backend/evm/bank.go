package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"strings"
)

func (m *EvmBackend) GetBankDenomsMetadata(coins sdk.Coins) map[string]banktypes.Metadata {
	denomsMetadata := make(map[string]banktypes.Metadata)
	for _, coin := range coins {
		res, err := m.queryClient.BankQueryClient.DenomMetadata(m.ctx, &banktypes.QueryDenomMetadataRequest{
			Denom: coin.Denom,
		})
		if err != nil || res == nil || coin.Denom == "" {
			continue
		}
		denomsMetadata[coin.Denom] = res.Metadata
	}

	if len(denomsMetadata) == 0 && len(coins) > 0 {
		// trying to insert denom metadata for the default RollApp coin

		for _, coin := range coins {
			if len(coin.Denom) < 2 {
				continue
			}

			prefixA := strings.HasPrefix(coin.Denom, "a")
			prefixU := strings.HasPrefix(coin.Denom, "u")
			if !prefixA && !prefixU {
				continue
			}

			// add pseudo data based on naming convention
			display := strings.ToUpper(coin.Denom[1:])
			denomsMetadata[coin.Denom] = banktypes.Metadata{
				DenomUnits: []*banktypes.DenomUnit{{
					Denom:    coin.Denom,
					Exponent: 0,
				}, {
					Denom: display,
					Exponent: func() uint32 {
						if prefixA {
							return 18
						}
						return 6
					}(),
				}},
				Base:    coin.Denom,
				Display: display,
				Name:    display,
				Symbol:  display,
			}
		}
	}

	return denomsMetadata
}
