package keeper

import (
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetChannel set a specific channel in the store from its index
func (k Keeper) SetChannel(ctx sdk.Context, channel types.Channel) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChannelKeyPrefix))
	b := k.cdc.MustMarshal(&channel)
	store.Set(types.ChannelKey(
		channel.Index,
	), b)
}

// GetChannel returns a channel from its index
func (k Keeper) GetChannel(
	ctx sdk.Context,
	index string,

) (val types.Channel, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChannelKeyPrefix))

	b := store.Get(types.ChannelKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveChannel removes a channel from the store
func (k Keeper) RemoveChannel(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChannelKeyPrefix))
	store.Delete(types.ChannelKey(
		index,
	))
}

// GetAllChannel returns all channel
func (k Keeper) GetAllChannel(ctx sdk.Context) (list []types.Channel) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChannelKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Channel
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
