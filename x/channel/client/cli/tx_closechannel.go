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

func CmdClosechannel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "closechannel [multisig-addr] [part-a] [coin-a] [part-b] [coin-b] [channelid]",
		Short: "Broadcast message closechannel",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argMultisigAddr := args[0]
			argPartA := args[1]
			argCoinA := args[2]
			argPartB := args[3]
			argCoinB := args[4]
			argChannelid := args[5]

			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(argMultisigAddr)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
			}

			_, err = sdk.AccAddressFromBech32(argPartA)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid toA address (%s)", err)
			}

			_, err = sdk.AccAddressFromBech32(argPartB)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid toB address (%s)", err)
			}

			var coinA, coinB []*sdk.Coin

			argcoin := strings.Split(argCoinA, ":")
			coinA = make([]*sdk.Coin, len(argcoin))

			for i, coin := range argcoin {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				coinA[i] = &c
			}

			argcoin = strings.Split(argCoinB, ":")
			coinB = make([]*sdk.Coin, len(argcoin))

			for i, coin := range argcoin {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				coinB[i] = &c
			}

			msg := types.NewMsgClosechannel(
				clientCtx.GetFromAddress().String(),
				argMultisigAddr,
				argPartA,
				coinA,
				argPartB,
				coinB,
				argChannelid,
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
