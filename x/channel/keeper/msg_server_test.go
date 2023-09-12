package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dungtt-astra/astra/testutil/keeper"
	"github.com/dungtt-astra/astra/x/channel/keeper"
	"github.com/dungtt-astra/astra/x/channel/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ChannelKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
