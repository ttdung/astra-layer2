package keeper

import (
	"context"
	"fmt"

	"channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Openchannel(goCtx context.Context, msg *types.MsgOpenchannel) (*types.MsgOpenchannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	addrA, err := sdk.AccAddressFromBech32(msg.PartA)
	if err != nil {
		return nil, err
	}

	addrB, err := sdk.AccAddressFromBech32(msg.PartB)
	if err != nil {
		return nil, err
	}

	multiAddr := msg.GetSigners()[0]

	if msg.CoinA.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, addrA, multiAddr, sdk.Coins{*msg.CoinA})
		if err != nil {
			return nil, err
		}
	}

	if msg.CoinB.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, addrB, multiAddr, sdk.Coins{*msg.CoinB})
		if err != nil {
			return nil, err
		}
	}

	indexStr := fmt.Sprintf("%s:%s:%s", msg.MultisigAddr, msg.CoinA.Denom, msg.Sequence)

	channel := types.Channel{
		Index:        indexStr,
		MultisigAddr: msg.MultisigAddr,
		PartA:        msg.PartA,
		PartB:        msg.PartB,
		Denom:        msg.CoinA.Denom,
	}

	k.Keeper.SetChannel(ctx, channel)

	return &types.MsgOpenchannelResponse{
		Id: indexStr,
	}, nil
}
