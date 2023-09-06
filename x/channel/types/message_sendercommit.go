package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSendercommit = "sendercommit"

var _ sdk.Msg = &MsgSendercommit{}

func NewMsgSendercommit(creator string, senderAddr string, channelid string, cointosender *sdk.Coin, cointohtlc *sdk.Coin, cointransfer *sdk.Coin, hashcodehtlc string, timelockhtlc uint64, hashcodedest string, timelockreceiver uint64, timelocksender uint64, multisigAddr string) *MsgSendercommit {
	return &MsgSendercommit{
		Creator:          creator,
		SenderAddr:       senderAddr,
		Channelid:        channelid,
		Cointosender:     cointosender,
		Cointohtlc:       cointohtlc,
		Cointransfer:     cointransfer,
		Hashcodehtlc:     hashcodehtlc,
		Timelockhtlc:     timelockhtlc,
		Hashcodedest:     hashcodedest,
		Timelockreceiver: timelockreceiver,
		Timelocksender:   timelocksender,
		MultisigAddr:     multisigAddr,
	}
}

func (msg *MsgSendercommit) Route() string {
	return RouterKey
}

func (msg *MsgSendercommit) Type() string {
	return TypeMsgSendercommit
}

func (msg *MsgSendercommit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSendercommit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendercommit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
