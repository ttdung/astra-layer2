package keeper

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Fund(goCtx context.Context, msg *types.MsgFund) (*types.MsgFundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_, err := sdk.AccAddressFromBech32(msg.MultisigAddr)
	if err != nil {
		return nil, err
	}

	val, found := k.Keeper.GetChannel(ctx, msg.Channelid)
	if !found {
		return nil, fmt.Errorf("ChannelID %v is not existing", msg.Channelid)
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

	multisigAddr, err := sdk.AccAddressFromBech32(val.MultisigAddr)
	if err != nil {
		return nil, err
	}

	cointoPartner := msg.CointoPartner

	coinChannel := k.Keeper.bankKeeper.GetBalance(ctx, multisigAddr, cointoPartner.Denom)

	ctx.Logger().Info("@@@@ balance of addr", val.MultisigAddr,
		" balance:", coinChannel.Amount.Uint64(), " cointoPartner:", cointoPartner.Amount.Uint64())

	// Send to LockTx (other) or HashTx (creator)
	if cointoPartner.Amount.IsPositive() {
		err = k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, multisigAddr, types.ModuleName, sdk.Coins{*cointoPartner})
		if err != nil {
			return nil, fmt.Errorf("@@@ SendCoinsFromAccountToModule failed balance of addr: %v, balance: %v", val.MultisigAddr, cointoPartner.Amount.Uint64())
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
		Cointohtlc:    cointoPartner,
		Channelid:     msg.Channelid,
	}

	k.Keeper.SetCommitment(ctx, commitment)

	if creatorAddr != msg.Creatoraddr {
		return nil, fmt.Errorf("not matching receiver address! expected: %v", creatorAddr)
	}

	// Send coin to creator of the funding commitment
	to, err := sdk.AccAddressFromBech32(msg.Creatoraddr)
	if err != nil {
		return nil, err
	}

	coin_fundside := coinChannel.Sub(*cointoPartner)
	if coin_fundside.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, multisigAddr, to, sdk.Coins{sdk.Coin{coin_fundside.Denom, coin_fundside.Amount}})
		if err != nil {
			return nil, fmt.Errorf("SendCoins failed balance of addr %v, balance %v, amount %v",
				val.MultisigAddr,
				coinChannel.Amount.Uint64(),
				coin_fundside.Amount.Uint64())
		}
	}

	k.Keeper.RemoveChannel(ctx, msg.Channelid)

	return &types.MsgFundResponse{
		Id: indexStr,
	}, nil

}
