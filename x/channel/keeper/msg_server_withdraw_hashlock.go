package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
)

func (k msgServer) WithdrawHashlock(goCtx context.Context, msg *types.MsgWithdrawHashlock) (*types.MsgWithdrawHashlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	val, found := k.Keeper.GetCommitment(ctx, msg.Index)
	if !found {
		return nil, fmt.Errorf("Commitment %v is not existing", msg.Index)
	}

	if val.Creatoraddr != msg.To {
		return nil, fmt.Errorf("Not matching receiver address! expected: %v", val.Creatoraddr)
	}

	hash := sha256.Sum256([]byte(msg.Secret))
	if val.Hashcode != base64.StdEncoding.EncodeToString(hash[:]) {
		return nil, fmt.Errorf("Wrong hash !")
	}

	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	if k.bankKeeper.BlockedAddr(to) {
		err = fmt.Errorf("%s is not allowed to receive funds", msg.To)
	} else {
		for _, coin := range val.Cointohtlc {
			if coin.Amount.IsPositive() {
				err = k.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, sdk.Coins{*coin})
				if err != nil {
					return nil, err
				}
			}
		}
	}

	k.Keeper.RemoveCommitment(ctx, msg.Index)

	return &types.MsgWithdrawHashlockResponse{}, nil
}
