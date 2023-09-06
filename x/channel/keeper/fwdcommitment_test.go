package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/AstraProtocol/astra/channel/testutil/keeper"
	"github.com/AstraProtocol/astra/channel/testutil/nullify"
	"github.com/AstraProtocol/astra/channel/x/channel/keeper"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNFwdcommitment(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Fwdcommitment {
	items := make([]types.Fwdcommitment, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetFwdcommitment(ctx, items[i])
	}
	return items
}

func TestFwdcommitmentGet(t *testing.T) {
	keeper, ctx := keepertest.ChannelKeeper(t)
	items := createNFwdcommitment(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetFwdcommitment(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestFwdcommitmentRemove(t *testing.T) {
	keeper, ctx := keepertest.ChannelKeeper(t)
	items := createNFwdcommitment(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveFwdcommitment(ctx,
			item.Index,
		)
		_, found := keeper.GetFwdcommitment(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestFwdcommitmentGetAll(t *testing.T) {
	keeper, ctx := keepertest.ChannelKeeper(t)
	items := createNFwdcommitment(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllFwdcommitment(ctx)),
	)
}
