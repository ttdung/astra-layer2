package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
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

			decCoin, err := sdk.ParseDecCoin(argCoinA)
			if err != nil {
				return err
			}
			coinA, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			decCoin, err = sdk.ParseDecCoin(argCoinB)
			if err != nil {
				return err
			}
			coinB, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			msg := types.NewMsgClosechannel(
				clientCtx.GetFromAddress().String(),
				argMultisigAddr,
				argPartA,
				&coinA,
				argPartB,
				&coinB,
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
