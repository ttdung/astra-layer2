package keeper

import (
	"context"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FwdcommitmentAll(goCtx context.Context, req *types.QueryAllFwdcommitmentRequest) (*types.QueryAllFwdcommitmentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var fwdcommitments []types.Fwdcommitment
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	fwdcommitmentStore := prefix.NewStore(store, types.KeyPrefix(types.FwdcommitmentKeyPrefix))

	pageRes, err := query.Paginate(fwdcommitmentStore, req.Pagination, func(key []byte, value []byte) error {
		var fwdcommitment types.Fwdcommitment
		if err := k.cdc.Unmarshal(value, &fwdcommitment); err != nil {
			return err
		}

		fwdcommitments = append(fwdcommitments, fwdcommitment)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllFwdcommitmentResponse{Fwdcommitment: fwdcommitments, Pagination: pageRes}, nil
}

func (k Keeper) Fwdcommitment(goCtx context.Context, req *types.QueryGetFwdcommitmentRequest) (*types.QueryGetFwdcommitmentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetFwdcommitment(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetFwdcommitmentResponse{Fwdcommitment: val}, nil
}
