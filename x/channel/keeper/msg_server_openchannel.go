package keeper

import (
	"context"
	"fmt"
	"github.com/dungtt-astra/astra/v3/x/channel/pubkey"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
)

func (k msgServer) Openchannel(goCtx context.Context, msg *types.MsgOpenchannel) (*types.MsgOpenchannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pubkeyA, err := pubkey.NewPKAccount(msg.PartA)
	//addrA, err := sdk.AccAddressFromBech32(msg.PartA)
	if err != nil {
		return nil, err
	}
	addrA := pubkeyA.AccAddress()

	pubkeyB, err := pubkey.NewPKAccount(msg.PartB)
	//addrB, err := sdk.AccAddressFromBech32(msg.PartB)
	if err != nil {
		return nil, err
	}
	addrB := pubkeyB.AccAddress()

	multiAddr := msg.GetSigners()[0]

	if msg.CoinA.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, addrA, multiAddr, sdk.Coins{*msg.CoinA})
		if err != nil {
			log.Println("-------------.. Err:", err.Error())
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
		Sequence:     msg.Sequence,
	}

	k.Keeper.SetChannel(ctx, channel)

	return &types.MsgOpenchannelResponse{
		Id: indexStr,
	}, nil
}
