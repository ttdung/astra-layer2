package keeper

import (
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetFwdcommitment set a specific fwdcommitment in the store from its index
func (k Keeper) SetFwdcommitment(ctx sdk.Context, fwdcommitment types.Fwdcommitment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FwdcommitmentKeyPrefix))
	b := k.cdc.MustMarshal(&fwdcommitment)
	store.Set(types.FwdcommitmentKey(
		fwdcommitment.Index,
	), b)
}

// GetFwdcommitment returns a fwdcommitment from its index
func (k Keeper) GetFwdcommitment(
	ctx sdk.Context,
	index string,

) (val types.Fwdcommitment, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FwdcommitmentKeyPrefix))

	b := store.Get(types.FwdcommitmentKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveFwdcommitment removes a fwdcommitment from the store
func (k Keeper) RemoveFwdcommitment(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FwdcommitmentKeyPrefix))
	store.Delete(types.FwdcommitmentKey(
		index,
	))
}

// GetAllFwdcommitment returns all fwdcommitment
func (k Keeper) GetAllFwdcommitment(ctx sdk.Context) (list []types.Fwdcommitment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FwdcommitmentKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Fwdcommitment
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
