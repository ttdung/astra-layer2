package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgReceiverwithdraw = "receiverwithdraw"

var _ sdk.Msg = &MsgReceiverwithdraw{}

func NewMsgReceiverwithdraw(creator string, transferindex string, to string, secret string) *MsgReceiverwithdraw {
	return &MsgReceiverwithdraw{
		Creator:       creator,
		Transferindex: transferindex,
		To:            to,
		Secret:        secret,
	}
}

func (msg *MsgReceiverwithdraw) Route() string {
	return RouterKey
}

func (msg *MsgReceiverwithdraw) Type() string {
	return TypeMsgReceiverwithdraw
}

func (msg *MsgReceiverwithdraw) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgReceiverwithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReceiverwithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
