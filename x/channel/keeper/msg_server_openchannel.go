package keeper

import (
	"context"
	"fmt"
	"github.com/dungtt-astra/astra/v3/x/channel/pubkey"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
)

func (k msgServer) Openchannel(goCtx context.Context, msg *types.MsgOpenchannel) (*types.MsgOpenchannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	pubkeyA, err := pubkey.NewPKAccount(msg.PartA)
	if err != nil {
		return nil, err
	}

	addrA := pubkeyA.AccAddress()

	pubkeyB, err := pubkey.NewPKAccount(msg.PartB)
	if err != nil {
		return nil, err
	}

	addrB := pubkeyB.AccAddress()

	multiAddr := msg.GetSigners()[0]

	// Verify multisig addr of the signer
	if strings.Compare(multiAddr.String(), msg.MultisigAddr) != 0 {
		return nil, fmt.Errorf("Wrong multisig address, expected:", msg.MultisigAddr)
	}

	for _, coin := range msg.CoinA {
		if coin.Amount.IsPositive() {
			err = k.bankKeeper.SendCoins(ctx, addrA, multiAddr, sdk.Coins{*coin})
			if err != nil {
				return nil, err
			}
		}
	}

	for _, coin := range msg.CoinB {
		if coin.Amount.IsPositive() {
			err = k.bankKeeper.SendCoins(ctx, addrB, multiAddr, sdk.Coins{*coin})
			if err != nil {
				return nil, err
			}
		}
	}

	indexStr := fmt.Sprintf("%s:%s:%s", msg.MultisigAddr, msg.CoinA[0].Denom, msg.Sequence)

	channel := types.Channel{
		Index:        indexStr,
		MultisigAddr: msg.MultisigAddr,
		PartA:        msg.PartA,
		PartB:        msg.PartB,
		Denom:        msg.CoinA[0].Denom,
		Sequence:     msg.Sequence,
	}

	k.Keeper.SetChannel(ctx, channel)

	return &types.MsgOpenchannelResponse{
		Id: indexStr,
	}, nil
}
