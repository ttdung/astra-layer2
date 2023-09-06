package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Senderwithdrawhashlock(goCtx context.Context, msg *types.MsgSenderwithdrawhashlock) (*types.MsgSenderwithdrawhashlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	val, found := k.Keeper.GetFwdcommitment(ctx, msg.Transferindex)
	if !found {
		return nil, fmt.Errorf("Fwdcommitment %v does not exist", msg.Transferindex)
	}

	if val.SenderAddr != msg.To {
		return nil, fmt.Errorf("not matching receiver address! expected: %v", val.SenderAddr)
	}

	hash := sha256.Sum256([]byte(msg.Secret))
	if val.Hashcodehtlc != base64.StdEncoding.EncodeToString(hash[:]) {
		return nil, fmt.Errorf("Wrong hash !")
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

	return &types.MsgSenderwithdrawhashlockResponse{}, nil
}
