package ante

import (
	"fmt"
	//"github.com/AstraProtocol/astra/v3/x/channel/types"
	"math/big"

	channelkeeper "github.com/AstraProtocol/astra/channel/x/channel/keeper"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
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
	params := lnmsg.ChannelKeeper.GetParams(ctx)
	fmt.Println("params=========================================:", params)
	//ethCfg := params.ChainConfig.EthereumConfig(gwd.evmKeeper.ChainID())
	blockHeight := big.NewInt(ctx.BlockHeight())
	//isLondon := ethCfg.IsLondon(blockHeight)

	fmt.Println("blockHeight=========================================:", blockHeight)

	authTx, ok := tx.(authsigning.SigVerifiableTx) //(sdk.FeeTx)
	if !ok {
		return next(ctx, tx, simulate)
	}

	msg := authTx.GetMsgs()[0]

	switch m := msg.(type) {
	case *types.MsgClosechannel:
		err = lnmsg.validateCloseChannelTx(ctx, authTx, m)
	case *types.MsgOpenchannel:
		err = lnmsg.validateOpenChannelTx(ctx, authTx, m)

	}

	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (lnmsg LnMsgDecorator) validateCloseChannelTx(ctx sdk.Context, authTx authsigning.SigVerifiableTx, m *types.MsgClosechannel) error {

	c, found := lnmsg.ChannelKeeper.GetChannel(ctx, m.Channelid)
	if !found {
		return fmt.Errorf("do not found channel id:", m.Channelid)
	}

	// verify correct multisig address or not
	if m.MultisigAddr != c.MultisigAddr {
		return fmt.Errorf("wrong multisig address, expected:", c.MultisigAddr)
	}

	multisig, err := sdk.AccAddressFromBech32(m.MultisigAddr)
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

	return nil
}
