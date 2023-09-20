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

func CmdReceivercommit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receivercommit [receiver-addr] [channelid] [cointoreceiver] [cointohtlc] [cointransfer] [hashcodehtlc] [timelockhtlc] [hashcodedest] [timelocksender] [multisig-addr]",
		Short: "Broadcast message receivercommit",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argReceiverAddr := args[0]
			argChannelid := args[1]
			//argCointoreceiver := args[2]
			//argCointohtlc := args[3]
			//argCointransfer := args[4]
			argHashcodehtlc := args[5]
			//argTimelockhtlc := args[6]
			argHashcodedest := args[7]
			//argTimelocksender := args[8]
			argMultisigAddr := args[9]

			cmd.Flags().Set(flags.FlagFrom, argMultisigAddr)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var argCointoreceiver, argCointohtlc, argCointransfer []*sdk.Coin

			arg2 := strings.Split(args[2], ":")
			argCointoreceiver = make([]*sdk.Coin, len(arg2))
			for i, coin := range arg2 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				argCointoreceiver[i] = &c
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

			argTimelocksender, err := strconv.ParseUint(args[8], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgReceivercommit(
				clientCtx.GetFromAddress().String(),
				argReceiverAddr,
				argChannelid,
				argCointoreceiver,
				argCointohtlc,
				argCointransfer,
				argHashcodehtlc,
				argTimelockhtlc,
				argHashcodedest,
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
