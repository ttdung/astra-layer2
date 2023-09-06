package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawHashlock = "withdraw_hashlock"

var _ sdk.Msg = &MsgWithdrawHashlock{}

func NewMsgWithdrawHashlock(creator string, to string, index string, secret string) *MsgWithdrawHashlock {
	return &MsgWithdrawHashlock{
		Creator: creator,
		To:      to,
		Index:   index,
		Secret:  secret,
	}
}

func (msg *MsgWithdrawHashlock) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawHashlock) Type() string {
	return TypeMsgWithdrawHashlock
}

func (msg *MsgWithdrawHashlock) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawHashlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawHashlock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
