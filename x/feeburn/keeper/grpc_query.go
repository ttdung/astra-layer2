package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/v3/cmd/config"
	"github.com/dungtt-astra/astra/v3/x/feeburn/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) TotalFeeBurn(c context.Context, request *types.QueryTotalFeeBurnRequest) (*types.QueryTotalFeeBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	totalFeeBurn := k.GetTotalFeeBurn(ctx)
	return &types.QueryTotalFeeBurnResponse{TotalFeeBurn: sdk.NewDecCoinFromDec(config.BaseDenom, totalFeeBurn)}, nil
}
