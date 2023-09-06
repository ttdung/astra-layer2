package keeper

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Receivercommit(goCtx context.Context, msg *types.MsgReceivercommit) (*types.MsgReceivercommitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	val, found := k.Keeper.GetChannel(ctx, msg.Channelid)
	if !found {
		return nil, fmt.Errorf("ChannelID %v is not existing", msg.Channelid)
	}

	var partnerAddr, creatorAddr string
	if msg.ReceiverAddr == val.PartA {
		partnerAddr = val.PartB
		creatorAddr = val.PartA
	} else {
		partnerAddr = val.PartA
		creatorAddr = val.PartB
	}

	toReceiver, err := sdk.AccAddressFromBech32(msg.ReceiverAddr)
	if err != nil {
		return nil, err
	}

	from, err := sdk.AccAddressFromBech32(msg.MultisigAddr)
	if err != nil {
		return nil, err
	}

	// Send coin to creator of commitment
	if msg.Cointoreceiver.Amount.IsPositive() {
		err = k.bankKeeper.SendCoins(ctx, from, toReceiver, sdk.Coins{*msg.Cointoreceiver})
		if err != nil {
			return nil, err
		}
	}

	// Send to HTLC
	htlcIndex := fmt.Sprintf("%s:%s", msg.Channelid, msg.Hashcodehtlc)
	Cointohtlc := msg.Cointohtlc
	if Cointohtlc.Amount.IsPositive() {
		err = k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.Coins{*Cointohtlc})
		if err != nil {
			return nil, err
		}
	}

	unlockBlockHtlc := msg.Timelockhtlc + uint64(ctx.BlockHeight())

	commitment := types.Commitment{
		Index:         htlcIndex,
		MultisigAddr:  msg.MultisigAddr,
		Creatoraddr:   creatorAddr,
		Partneraddr:   partnerAddr,
		Hashcode:      msg.Hashcodehtlc,
		Numblock:      unlockBlockHtlc,
		Cointocreator: msg.Cointoreceiver,
		Cointohtlc:    msg.Cointohtlc,
		Channelid:     msg.Channelid,
	}

	k.Keeper.SetCommitment(ctx, commitment)

	// Send to FwdContract
	CointoFC := msg.Cointransfer
	if CointoFC.Amount.IsPositive() {
		err = k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.Coins{*CointoFC})
		if err != nil {
			return nil, err
		}
	}

	creator := "receiver"
	transferIndex := fmt.Sprintf("%s:%s:%s", msg.Channelid, msg.Hashcodedest, creator)

	Timelocksender := msg.Timelocksender

	fwscommitment := types.Fwdcommitment{
		Index:            transferIndex,
		Channelid:        msg.Channelid,
		MultisigAddr:     msg.MultisigAddr,
		SenderAddr:       partnerAddr,
		ReceiverAddr:     creatorAddr, // this commitment is created by receiver side
		Timelockreceiver: 0,
		Timelocksender:   Timelocksender,
		Hashcodehtlc:     msg.Hashcodehtlc,
		Hashcodedest:     msg.Hashcodedest,
		Cointransfer:     msg.Cointransfer,
		Creator:          creator,
	}

	k.Keeper.SetFwdcommitment(ctx, fwscommitment)

	k.Keeper.RemoveChannel(ctx, msg.Channelid)

	return &types.MsgReceivercommitResponse{
		HtlcIndex:     htlcIndex,
		TransferIndex: transferIndex,
	}, nil
}
