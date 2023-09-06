package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
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

			decCoin, err := sdk.ParseDecCoin(args[2])
			if err != nil {
				return err
			}
			argCointoreceiver, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			decCoin, err = sdk.ParseDecCoin(args[3])
			if err != nil {
				return err
			}
			argCointohtlc, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			decCoin, err = sdk.ParseDecCoin(args[4])
			if err != nil {
				return err
			}
			argCointransfer, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

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
				&argCointoreceiver,
				&argCointohtlc,
				&argCointransfer,
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
