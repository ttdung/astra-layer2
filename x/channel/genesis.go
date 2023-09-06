package channel

import (
	"github.com/AstraProtocol/astra/channel/x/channel/keeper"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the channel
	for _, elem := range genState.ChannelList {
		k.SetChannel(ctx, elem)
	}
	// Set all the commitment
	for _, elem := range genState.CommitmentList {
		k.SetCommitment(ctx, elem)
	}
	// Set all the fwdcommitment
	for _, elem := range genState.FwdcommitmentList {
		k.SetFwdcommitment(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ChannelList = k.GetAllChannel(ctx)
	genesis.CommitmentList = k.GetAllCommitment(ctx)
	genesis.FwdcommitmentList = k.GetAllFwdcommitment(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
