package keeper

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Acceptfund(goCtx context.Context, msg *types.MsgAcceptfund) (*types.MsgAcceptfundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_, err := sdk.AccAddressFromBech32(msg.MultisigAddr)
	if err != nil {
		return nil, err
	}

	val, found := k.Keeper.GetChannel(ctx, msg.Channelid)
	if !found {
		return nil, fmt.Errorf("ChannelID is not existing")
	}

	if msg.MultisigAddr != val.MultisigAddr {
		return nil, fmt.Errorf("Not matching multisig address!")
	}

	if msg.Creatoraddr != val.PartA && msg.Creatoraddr != val.PartB {
		return nil, fmt.Errorf("Not matching any part in this channel!")
	}

	var partnerAddr, creatorAddr string
	if msg.Creatoraddr == val.PartA {
		partnerAddr = val.PartB
		creatorAddr = val.PartA
	} else {
		partnerAddr = val.PartA
		creatorAddr = val.PartB
	}

	from, err := sdk.AccAddressFromBech32(val.MultisigAddr)
	if err != nil {
		return nil, err
	}

	coin_acceptside := msg.CointoCreator
	coin_channel := k.Keeper.bankKeeper.GetBalance(ctx, from, coin_acceptside.Denom)

	// Send coin to accepted side first
	to, err := sdk.AccAddressFromBech32(msg.Creatoraddr)
	if err != nil {
		return nil, err
	}

	if coin_acceptside.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, from, to, sdk.Coins{*coin_acceptside})
		if err != nil {
			return nil, fmt.Errorf("SendCoins failed balance of addr %v, balance: %v, required amt: %v",
				val.MultisigAddr,
				coin_channel.Amount.Uint64(),
				coin_acceptside.Amount.Uint64())
		}
	}

	// Send the remain coin to HTLC = coin_channel - coin_acceptside
	coin_htlc := coin_channel.Sub(*coin_acceptside)

	if coin_htlc.Amount.IsPositive() {
		err = k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.Coins{coin_htlc})
		if err != nil {
			return nil, fmt.Errorf("@@@ SendCoinsFromAccountToModule failed, Addr: %v, balance: %v ",
				val.MultisigAddr, coin_channel.Amount.Uint64())
		}
	}

	indexStr := fmt.Sprintf("%s:%s", msg.Channelid, msg.Hashcode)

	unlockBlock := msg.Numblock + uint64(ctx.BlockHeight())

	commitment := types.Commitment{
		Index:         indexStr,
		MultisigAddr:  msg.MultisigAddr,
		Creatoraddr:   creatorAddr,
		Partneraddr:   partnerAddr,
		Hashcode:      msg.Hashcode,
		Numblock:      unlockBlock,
		Cointocreator: nil,
		Cointohtlc:    &coin_htlc,
		Channelid:     msg.Channelid,
	}

	k.Keeper.SetCommitment(ctx, commitment)

	if creatorAddr != msg.Creatoraddr {
		return nil, fmt.Errorf("Not matching receiver address! expected: %v", creatorAddr)
	}

	k.Keeper.RemoveChannel(ctx, msg.Channelid)

	return &types.MsgAcceptfundResponse{
		Id: indexStr,
	}, nil
}
