package keeper

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/x/channel/types"
)

func (k msgServer) Openchannel(goCtx context.Context, msg *types.MsgOpenchannel) (*types.MsgOpenchannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	log.Println("========... OPEN CHANNEL:", msg)

	addrA, err := sdk.AccAddressFromBech32(msg.PartA)
	if err != nil {
		return nil, err
	}
	log.Println("========... OPEN CHANNEL: step 1 addrA:", addrA)
	addrB, err := sdk.AccAddressFromBech32(msg.PartB)
	if err != nil {
		return nil, err
	}

	log.Println("========... OPEN CHANNEL: step 2 addrB:", addrB)
	multiAddr := msg.GetSigners()[0]

	log.Println("========... OPEN CHANNEL: Multiaddr:", multiAddr)

	if msg.CoinA.Amount.IsPositive() {
		log.Println("========... OPEN CHANNEL: step 3 ")
		err = k.bankKeeper.SendCoins(ctx, addrA, multiAddr, sdk.Coins{*msg.CoinA})
		if err != nil {
			log.Println("========... OPEN CHANNEL: err1 ", err.Error())
			return nil, err
		}
	}

	if msg.CoinB.Amount.IsPositive() {
		log.Println("========... OPEN CHANNEL: step 4 ")
		err = k.bankKeeper.SendCoins(ctx, addrB, multiAddr, sdk.Coins{*msg.CoinB})
		if err != nil {
			log.Println("========... OPEN CHANNEL: err2 ", err.Error())
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
