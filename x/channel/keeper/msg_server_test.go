package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dungtt-astra/astra/v3/testutil/keeper"
	"github.com/dungtt-astra/astra/v3/x/channel/keeper"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ChannelKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
