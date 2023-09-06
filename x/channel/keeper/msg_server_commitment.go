package keeper

import (
	"context"
	"fmt"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Commitment(goCtx context.Context, msg *types.MsgCommitment) (*types.MsgCommitmentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	from, err := sdk.AccAddressFromBech32(msg.MultisigAddr)
	if err != nil {
		return nil, err
	}

	toCreator, err := sdk.AccAddressFromBech32(msg.Creatoraddr)
	if err != nil {
		return nil, err
	}

	_, err = sdk.AccAddressFromBech32(msg.Partneraddr)
	if err != nil {
		return nil, err
	}

	// Send coin to creator of commitment
	if msg.Cointocreator.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, from, toCreator, sdk.Coins{*msg.Cointocreator})
		if err != nil {
			return nil, err
		}
	}

	// Send to HTLC
	indexStr := fmt.Sprintf("%s:%s", msg.Channelid, msg.Hashcode)
	if msg.Cointohtlc.Amount.IsPositive() {
		err = k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.Coins{
			sdk.Coin{
				Denom:  msg.Cointohtlc.Denom,
				Amount: msg.Cointohtlc.Amount,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	unlockBlock := msg.Numblock + uint64(ctx.BlockHeight())

	commitment := types.Commitment{
		Index:         indexStr,
		MultisigAddr:  msg.MultisigAddr,
		Creatoraddr:   msg.Creatoraddr,
		Partneraddr:   msg.Partneraddr,
		Hashcode:      msg.Hashcode,
		Numblock:      unlockBlock,
		Cointocreator: msg.Cointocreator,
		Cointohtlc:    msg.Cointohtlc,
		Channelid:     msg.Channelid,
	}
	k.Keeper.SetCommitment(ctx, commitment)

	k.Keeper.RemoveChannel(ctx, msg.Channelid)

	return &types.MsgCommitmentResponse{
		Id: indexStr,
	}, nil
}
