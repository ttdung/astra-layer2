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

func CmdFund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund [creatoraddr] [channelid] [cointo-partner] [hashcode] [numblock] [multisig-addr]",
		Short: "Broadcast message fund",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argCreatoraddr := args[0]
			argChannelid := args[1]
			argHashcode := args[3]
			argMultisigAddr := args[5]

			cmd.Flags().Set(flags.FlagFrom, argMultisigAddr)
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			decCoin, err := sdk.ParseDecCoin(args[2])
			if err != nil {
				return err
			}
			argCointoPartner, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			argNumblock, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgFund(
				clientCtx.GetFromAddress().String(),
				argCreatoraddr,
				argChannelid,
				&argCointoPartner,
				argHashcode,
				argNumblock,
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
