package keeper

import (
	"context"
	"fmt"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
)

func (k msgServer) Openchannel(goCtx context.Context, msg *types.MsgOpenchannel) (*types.MsgOpenchannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	//kb := keyring.Keyring()

	addrA, err := sdk.AccAddressFromBech32(msg.PartA)
	if err != nil {
		return nil, err
	}

	addrB, err := sdk.AccAddressFromBech32(msg.PartB)
	if err != nil {
		return nil, err
	}

	multiAddr := msg.GetSigners()[0]

	// Verify multisig addr vs each single key
	if strings.Compare(multiAddr.String(), msg.MultisigAddr) != 0 {
		panic("Wrong multisig address")
	}

	for _, coin := range msg.CoinA {
		if coin.Amount.IsPositive() {
			err = k.bankKeeper.SendCoins(ctx, addrA, multiAddr, sdk.Coins{*coin})
			if err != nil {
				log.Println("-------------.. Err:", err.Error())
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
