package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSenderwithdrawtimelock = "senderwithdrawtimelock"

var _ sdk.Msg = &MsgSenderwithdrawtimelock{}

func NewMsgSenderwithdrawtimelock(creator string, transferindex string, to string) *MsgSenderwithdrawtimelock {
	return &MsgSenderwithdrawtimelock{
		Creator:       creator,
		Transferindex: transferindex,
		To:            to,
	}
}

func (msg *MsgSenderwithdrawtimelock) Route() string {
	return RouterKey
}

func (msg *MsgSenderwithdrawtimelock) Type() string {
	return TypeMsgSenderwithdrawtimelock
}

func (msg *MsgSenderwithdrawtimelock) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSenderwithdrawtimelock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSenderwithdrawtimelock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
