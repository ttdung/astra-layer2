package channel_test

import (
	"testing"

	keepertest "channel/testutil/keeper"
	"channel/testutil/nullify"
	"channel/x/channel"
	"channel/x/channel/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ChannelList: []types.Channel{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ChannelKeeper(t)
	channel.InitGenesis(ctx, *k, genesisState)
	got := channel.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ChannelList, got.ChannelList)
	// this line is used by starport scaffolding # genesis/test/assert
}
