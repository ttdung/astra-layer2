package keeper_test

import (
	"context"
	"testing"

	keepertest "channel/testutil/keeper"
	"channel/x/channel/keeper"
	"channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ChannelKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
