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

func CmdSenderwithdrawhashlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "senderwithdrawhashlock [transferindex] [to] [secret]",
		Short: "Broadcast message senderwithdrawhashlock",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTransferindex := args[0]
			argTo := args[1]
			argSecret := args[2]

			cmd.Flags().Set(flags.FlagFrom, args[1])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSenderwithdrawhashlock(
				clientCtx.GetFromAddress().String(),
				argTransferindex,
				argTo,
				argSecret,
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
