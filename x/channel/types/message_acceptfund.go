package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAcceptfund = "acceptfund"

var _ sdk.Msg = &MsgAcceptfund{}

func NewMsgAcceptfund(creator string, creatoraddr string, channelid string, cointoCreator *sdk.Coin, hashcode string, numblock uint64, multisigAddr string) *MsgAcceptfund {
	return &MsgAcceptfund{
		Creator:       creator,
		Creatoraddr:   creatoraddr,
		Channelid:     channelid,
		CointoCreator: cointoCreator,
		Hashcode:      hashcode,
		Numblock:      numblock,
		MultisigAddr:  multisigAddr,
	}
}

func (msg *MsgAcceptfund) Route() string {
	return RouterKey
}

func (msg *MsgAcceptfund) Type() string {
	return TypeMsgAcceptfund
}

func (msg *MsgAcceptfund) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAcceptfund) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAcceptfund) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
