package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSenderwithdrawhashlock = "senderwithdrawhashlock"

var _ sdk.Msg = &MsgSenderwithdrawhashlock{}

func NewMsgSenderwithdrawhashlock(creator string, transferindex string, to string, secret string) *MsgSenderwithdrawhashlock {
	return &MsgSenderwithdrawhashlock{
		Creator:       creator,
		Transferindex: transferindex,
		To:            to,
		Secret:        secret,
	}
}

func (msg *MsgSenderwithdrawhashlock) Route() string {
	return RouterKey
}

func (msg *MsgSenderwithdrawhashlock) Type() string {
	return TypeMsgSenderwithdrawhashlock
}

func (msg *MsgSenderwithdrawhashlock) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSenderwithdrawhashlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSenderwithdrawhashlock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
