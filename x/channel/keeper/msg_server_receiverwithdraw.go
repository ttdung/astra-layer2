package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Receiverwithdraw(goCtx context.Context, msg *types.MsgReceiverwithdraw) (*types.MsgReceiverwithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	val, found := k.Keeper.GetFwdcommitment(ctx, msg.Transferindex)
	if !found {
		return nil, fmt.Errorf("Transferindex %v is not existing", msg.Transferindex)
	}

	if val.ReceiverAddr != msg.To {
		return nil, fmt.Errorf("not matching receiver address! expected: %v", val.ReceiverAddr)
	}

	hash := sha256.Sum256([]byte(msg.Secret))
	if val.Creator == "receiver" {
		// Receiver created this Forward Contract =>
		// the withdrawer could provide either the destination secret or the sender secret
		if val.Hashcodedest != base64.StdEncoding.EncodeToString(hash[:]) &&
			val.Hashcodehtlc != base64.StdEncoding.EncodeToString(hash[:]) {
			return nil, fmt.Errorf("Wrong hash !")
		}
	} else {
		// sender created this Forward Contract =>
		// the withdrawer has to provide the destination secret
		if val.Hashcodedest != base64.StdEncoding.EncodeToString(hash[:]) {
			return nil, fmt.Errorf("Wrong hash !")
		}
	}

	if val.Timelockreceiver > uint64(ctx.BlockHeight()) {
		return nil, fmt.Errorf("wait until valid block height, expected height %v", ctx.BlockHeight())
	}

	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, sdk.Coins{*val.Cointransfer})
	if err != nil {
		return nil, err
	}

	k.Keeper.RemoveFwdcommitment(ctx, msg.Transferindex)

	return &types.MsgReceiverwithdrawResponse{}, nil
}
