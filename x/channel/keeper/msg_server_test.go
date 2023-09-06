package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/AstraProtocol/astra/channel/testutil/keeper"
	"github.com/AstraProtocol/astra/channel/x/channel/keeper"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ChannelKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
