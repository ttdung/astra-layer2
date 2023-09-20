package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgReceivercommit = "receivercommit"

var _ sdk.Msg = &MsgReceivercommit{}

func NewMsgReceivercommit(creator string, receiverAddr string, channelid string, cointoreceiver []*sdk.Coin, cointohtlc []*sdk.Coin, cointransfer []*sdk.Coin, hashcodehtlc string, timelockhtlc uint64, hashcodedest string, timelocksender uint64, multisigAddr string) *MsgReceivercommit {
	return &MsgReceivercommit{
		Creator:        creator,
		ReceiverAddr:   receiverAddr,
		Channelid:      channelid,
		Cointoreceiver: cointoreceiver,
		Cointohtlc:     cointohtlc,
		Cointransfer:   cointransfer,
		Hashcodehtlc:   hashcodehtlc,
		Timelockhtlc:   timelockhtlc,
		Hashcodedest:   hashcodedest,
		Timelocksender: timelocksender,
		MultisigAddr:   multisigAddr,
	}
}

func (msg *MsgReceivercommit) Route() string {
	return RouterKey
}

func (msg *MsgReceivercommit) Type() string {
	return TypeMsgReceivercommit
}

func (msg *MsgReceivercommit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgReceivercommit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReceivercommit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
