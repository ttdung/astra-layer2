package keeper

import (
	"context"

	"channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ChannelAll(goCtx context.Context, req *types.QueryAllChannelRequest) (*types.QueryAllChannelResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var channels []types.Channel
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	channelStore := prefix.NewStore(store, types.KeyPrefix(types.ChannelKeyPrefix))

	pageRes, err := query.Paginate(channelStore, req.Pagination, func(key []byte, value []byte) error {
		var channel types.Channel
		if err := k.cdc.Unmarshal(value, &channel); err != nil {
			return err
		}

		channels = append(channels, channel)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllChannelResponse{Channel: channels, Pagination: pageRes}, nil
}

func (k Keeper) Channel(goCtx context.Context, req *types.QueryGetChannelRequest) (*types.QueryGetChannelResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetChannel(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetChannelResponse{Channel: val}, nil
}
