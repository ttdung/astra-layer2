package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgOpenchannel = "openchannel"

var _ sdk.Msg = &MsgOpenchannel{}

func NewMsgOpenchannel(creator string, partA string, partB string, coinA []*sdk.Coin, coinB []*sdk.Coin, multisigAddr string, sequence string) *MsgOpenchannel {
	return &MsgOpenchannel{
		Creator:      creator,
		PartA:        partA,
		PartB:        partB,
		CoinA:        coinA,
		CoinB:        coinB,
		MultisigAddr: multisigAddr,
		Sequence:     sequence,
	}
}

func (msg *MsgOpenchannel) Route() string {
	return RouterKey
}

func (msg *MsgOpenchannel) Type() string {
	return TypeMsgOpenchannel
}

func (msg *MsgOpenchannel) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgOpenchannel) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpenchannel) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
