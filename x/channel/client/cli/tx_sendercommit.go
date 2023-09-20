package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dungtt-astra/astra/v3/x/channel/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSendercommit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendercommit [sender-addr] [channelid] [cointosender] [cointohtlc] [cointransfer] [hashcodehtlc] [timelockhtlc] [hashcodedest] [timelockreceiver] [timelocksender] [multisig-addr]",
		Short: "Broadcast message sendercommit",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSenderAddr := args[0]
			argChannelid := args[1]
			//argCointosender := args[2]
			//argCointohtlc := args[3]
			//argCointransfer := args[4]
			argHashcodehtlc := args[5]
			//argTimelockhtlc := args[6]
			argHashcodedest := args[7]
			//argTimelockreceiver := args[8]
			//argTimelocksender := args[9]
			argMultisigAddr := args[10]

			cmd.Flags().Set(flags.FlagFrom, argMultisigAddr)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var argCointosender, argCointohtlc, argCointransfer []*sdk.Coin

			arg2 := strings.Split(args[2], ":")
			argCointosender = make([]*sdk.Coin, len(arg2))
			for i, coin := range arg2 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				argCointosender[i] = &c
			}

			arg3 := strings.Split(args[3], ":")
			argCointohtlc = make([]*sdk.Coin, len(arg3))
			for i, coin := range arg3 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				argCointohtlc[i] = &c
			}

			arg4 := strings.Split(args[4], ":")
			argCointransfer = make([]*sdk.Coin, len(arg4))
			for i, coin := range arg4 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				argCointransfer[i] = &c
			}

			argTimelockhtlc, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}

			argTimelockreceiver, err := strconv.ParseUint(args[8], 10, 64)
			if err != nil {
				return err
			}

			argTimelocksender, err := strconv.ParseUint(args[9], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendercommit(
				clientCtx.GetFromAddress().String(),
				argSenderAddr,
				argChannelid,
				argCointosender,
				argCointohtlc,
				argCointransfer,
				argHashcodehtlc,
				argTimelockhtlc,
				argHashcodedest,
				argTimelockreceiver,
				argTimelocksender,
				argMultisigAddr,
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
