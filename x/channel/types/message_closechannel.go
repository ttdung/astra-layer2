package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClosechannel = "closechannel"

var _ sdk.Msg = &MsgClosechannel{}

func NewMsgClosechannel(creator string, multisigAddr string, partA string, coinA *sdk.Coin, partB string, coinB *sdk.Coin, channelid string) *MsgClosechannel {
	return &MsgClosechannel{
		Creator:      creator,
		MultisigAddr: multisigAddr,
		PartA:        partA,
		CoinA:        coinA,
		PartB:        partB,
		CoinB:        coinB,
		Channelid:    channelid,
	}
}

func (msg *MsgClosechannel) Route() string {
	return RouterKey
}

func (msg *MsgClosechannel) Type() string {
	return TypeMsgClosechannel
}

func (msg *MsgClosechannel) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClosechannel) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClosechannel) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
