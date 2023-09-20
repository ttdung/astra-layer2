package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFund = "fund"

var _ sdk.Msg = &MsgFund{}

func NewMsgFund(creator string, creatoraddr string, channelid string, cointoPartner []*sdk.Coin, hashcode string, numblock uint64, multisigAddr string) *MsgFund {
	return &MsgFund{
		Creator:       creator,
		Creatoraddr:   creatoraddr,
		Channelid:     channelid,
		CointoPartner: cointoPartner,
		Hashcode:      hashcode,
		Numblock:      numblock,
		MultisigAddr:  multisigAddr,
	}
}

func (msg *MsgFund) Route() string {
	return RouterKey
}

func (msg *MsgFund) Type() string {
	return TypeMsgFund
}

func (msg *MsgFund) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFund) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFund) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
