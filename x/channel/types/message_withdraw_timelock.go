package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawTimelock = "withdraw_timelock"

var _ sdk.Msg = &MsgWithdrawTimelock{}

func NewMsgWithdrawTimelock(creator string, to string, index string) *MsgWithdrawTimelock {
	return &MsgWithdrawTimelock{
		Creator: creator,
		To:      to,
		Index:   index,
	}
}

func (msg *MsgWithdrawTimelock) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawTimelock) Type() string {
	return TypeMsgWithdrawTimelock
}

func (msg *MsgWithdrawTimelock) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawTimelock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawTimelock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
