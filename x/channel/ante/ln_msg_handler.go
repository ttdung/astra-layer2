package ante

import (
	"fmt"
	"github.com/dungtt-astra/astra/v3/x/channel/pubkey"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	channelkeeper "github.com/dungtt-astra/astra/v3/x/channel/keeper"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
)

// GasWantedDecorator keeps track of the gasWanted amount on the current block in transient store
// for BaseFee calculation.
// NOTE: This decorator does not perform any validation
type LnMsgDecorator struct {
	ChannelKeeper *channelkeeper.Keeper
	BankKeeper    evmtypes.BankKeeper
}

// NewGasWantedDecorator creates a new NewGasWantedDecorator
func NewLnMsgDecorator(channelkeeper *channelkeeper.Keeper,
	bankkeeper evmtypes.BankKeeper) LnMsgDecorator {
	return LnMsgDecorator{
		channelkeeper,
		bankkeeper,
	}
}

func (lnmsg LnMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	authTx, ok := tx.(authsigning.SigVerifiableTx) //(sdk.FeeTx)
	if !ok {
		return next(ctx, tx, simulate)
	}

	msg := authTx.GetMsgs()[0]
	fmt.Println("msg==============================================:", msg)

	switch m := msg.(type) {
	case *types.MsgClosechannel:
		err = lnmsg.validateCloseChannelTx(ctx, authTx, m)
	case *types.MsgOpenchannel:
		err = lnmsg.validateOpenChannelTx(ctx, authTx, m)
	case *types.MsgCommitment:
		err = lnmsg.validateCommitmentTx(ctx, authTx, m)
	case *types.MsgWithdrawHashlock:
		err = lnmsg.validateWithdrawHashlockTx(m)
	case *types.MsgWithdrawTimelock:
		err = lnmsg.validateWithdrawTimelockTx(m)
	case *types.MsgFund:
		err = lnmsg.validateFundTx(ctx, authTx, m)
	case *types.MsgAcceptfund:
		err = lnmsg.validateAcceptFundTx(ctx, authTx, m)
	}

	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (lnmsg LnMsgDecorator) isMatchingMultisig(ctx sdk.Context, channelId string, multisig string) error {
	c, found := lnmsg.ChannelKeeper.GetChannel(ctx, channelId)
	if !found {
		return fmt.Errorf("do not found channel id:", channelId)
	}

	// verify correct multisig address or not
	if multisig != c.MultisigAddr {
		return fmt.Errorf("wrong multisig address, expected:", c.MultisigAddr)
	}

	return nil
}

func (lnmsg LnMsgDecorator) validateWithdrawTimelockTx(m *types.MsgWithdrawTimelock) error {

	_, err := sdk.AccAddressFromBech32(m.To)

	return err
}

func (lnmsg LnMsgDecorator) validateWithdrawHashlockTx(m *types.MsgWithdrawHashlock) error {

	_, err := sdk.AccAddressFromBech32(m.To)

	return err
}

func (lnmsg LnMsgDecorator) validateAcceptFundTx(ctx sdk.Context, authTx authsigning.SigVerifiableTx, m *types.MsgAcceptfund) error {

	if err := lnmsg.isMatchingMultisig(ctx, m.Channelid, m.MultisigAddr); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(m.Creatoraddr)
	if err != nil {
		return err
	}

	multisig, err := sdk.AccAddressFromBech32(m.MultisigAddr)
	if err != nil {
		return err
	}

	// verify right signer or not
	if !multisig.Equals(authTx.GetSigners()[0]) {
		return fmt.Errorf("wrong signer, expected:", multisig.String())
	}

	amt := lnmsg.BankKeeper.GetBalance(ctx, multisig, m.CointoCreator.Denom)

	// verify amount to withdraw
	if m.CointoCreator.Amount.Int64() > amt.Amount.Int64() {
		return fmt.Errorf("exceed amount of token can be sent")
	}

	return nil
}

func (lnmsg LnMsgDecorator) validateFundTx(ctx sdk.Context, authTx authsigning.SigVerifiableTx, m *types.MsgFund) error {

	if err := lnmsg.isMatchingMultisig(ctx, m.Channelid, m.MultisigAddr); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(m.Creatoraddr)
	if err != nil {
		return err
	}

	multisig, err := sdk.AccAddressFromBech32(m.MultisigAddr)
	if err != nil {
		return err
	}

	// verify right signer or not
	if !multisig.Equals(authTx.GetSigners()[0]) {
		return fmt.Errorf("wrong signer, expected:", multisig.String())
	}

	amt := lnmsg.BankKeeper.GetBalance(ctx, multisig, m.CointoPartner.Denom)

	// verify amount to withdraw
	if m.CointoPartner.Amount.Int64() > amt.Amount.Int64() {
		return fmt.Errorf("exceed amount of token can be sent")
	}

	return nil
}

func (lnmsg LnMsgDecorator) validateCommitmentTx(ctx sdk.Context, authTx authsigning.SigVerifiableTx, m *types.MsgCommitment) error {

	if err := lnmsg.isMatchingMultisig(ctx, m.Channelid, m.MultisigAddr); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(m.Creatoraddr)
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(m.Partneraddr)
	if err != nil {
		return err
	}

	multisig, err := sdk.AccAddressFromBech32(m.MultisigAddr)
	if err != nil {
		return err
	}

	// verify right signer or not
	if !multisig.Equals(authTx.GetSigners()[0]) {
		return fmt.Errorf("wrong signer, expected:", multisig.String())
	}

	amt := lnmsg.BankKeeper.GetBalance(ctx, multisig, m.Cointocreator.Denom)

	// verify amount to withdraw
	if m.Cointohtlc.Amount.Int64()+m.Cointocreator.Amount.Int64() > amt.Amount.Int64() {
		return fmt.Errorf("exceed amount of token can be sent")
	}

	return nil
}

func (lnmsg LnMsgDecorator) validateCloseChannelTx(ctx sdk.Context, authTx authsigning.SigVerifiableTx, m *types.MsgClosechannel) error {

	if err := lnmsg.isMatchingMultisig(ctx, m.Channelid, m.MultisigAddr); err != nil {
		return err
	}

	multisig, err := sdk.AccAddressFromBech32(m.MultisigAddr)
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(m.PartA)
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(m.PartB)
	if err != nil {
		return err
	}

	// verify right signer or not
	if !multisig.Equals(authTx.GetSigners()[0]) {
		return fmt.Errorf("wrong signer, expected:", multisig.String())
	}

	amt := lnmsg.BankKeeper.GetBalance(ctx, multisig, m.CoinA.Denom)

	// verify amount to withdraw
	if m.CoinA.Amount.Int64()+m.CoinB.Amount.Int64() > amt.Amount.Int64() {
		return fmt.Errorf("exceed amount of token can be withdrawn", m.Channelid)
	}

	return nil
}

func (lnmsg LnMsgDecorator) validateOpenChannelTx(ctx sdk.Context, authTx authsigning.SigVerifiableTx, m *types.MsgOpenchannel) error {

	multisig, err := sdk.AccAddressFromBech32(m.MultisigAddr)
	if err != nil {
		return err
	}

	// verify right signer or not
	if !multisig.Equals(authTx.GetSigners()[0]) {
		return fmt.Errorf("wrong signer, expected:", multisig.String())
	}

	// validate the same coin
	if m.CoinA.Denom != m.CoinB.Denom {
		return fmt.Errorf("cannot create channel from different coin denoms")
	}

	// verify each party
	pubkeyA, err := pubkey.NewPKAccount(m.PartA)
	if err != nil {
		return err
	}
	pubkeyB, err := pubkey.NewPKAccount(m.PartB)
	if err != nil {
		return err
	}

	multisigAddr, _, err := pubkey.CreateMulSignAccountFromTwoAccount(pubkeyA.PublicKey(), pubkeyB.PublicKey(), 2)
	if err != nil {
		return err
	}

	//fmt.Println("====================...multisigAddr:", multisigAddr)
	//fmt.Println("====================...m.MultisigAddr:", m.MultisigAddr)

	if multisigAddr != m.MultisigAddr {
		return fmt.Errorf("Multisig and parties do not match")
	}

	return nil
}
