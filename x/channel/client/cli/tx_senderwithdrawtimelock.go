package cli

import (
	"strconv"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSenderwithdrawtimelock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "senderwithdrawtimelock [transferindex] [to]",
		Short: "Broadcast message senderwithdrawtimelock",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTransferindex := args[0]
			argTo := args[1]

			cmd.Flags().Set(flags.FlagFrom, args[1])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSenderwithdrawtimelock(
				clientCtx.GetFromAddress().String(),
				argTransferindex,
				argTo,
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
