package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdOpenchannel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openchannel [part-a] [part-b] [coin-a] [coin-b] [multisig-addr] [sequence]",
		Short: "Broadcast message openchannel",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPartA := args[0]
			argPartB := args[1]
			argMultisigAddr := args[4]
			argSequence := args[5]

			_, err = sdk.AccAddressFromBech32(argPartA)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid partA address (%s)", err)
			}

			_, err = sdk.AccAddressFromBech32(argPartB)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid partB address (%s)", err)
			}

			_, err = sdk.AccAddressFromBech32(argMultisigAddr)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid multisig address (%s)", err)
			}

			cmd.Flags().Set(flags.FlagFrom, argMultisigAddr)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var coinA, coinB []*sdk.Coin
			arg2 := strings.Split(args[2], ":")

			coinA = make([]*sdk.Coin, len(arg2))
			for i, coin := range arg2 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				coinA[i] = &c
			}

			arg3 := strings.Split(args[3], ":")
			coinB = make([]*sdk.Coin, len(arg3))
			for i, coin := range arg3 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				coinB[i] = &c
			}

			msg := types.NewMsgOpenchannel(
				clientCtx.GetFromAddress().String(),
				argPartA,
				argPartB,
				coinA,
				coinB,
				argMultisigAddr,
				argSequence,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
