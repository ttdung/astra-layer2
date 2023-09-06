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

func CmdCommitment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commitment [multisig-addr] [creatoraddr] [partneraddr] [hashcode] [numblock] [cointocreator] [cointohtlc] [channelid]",
		Short: "Broadcast message commitment",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argMultisigAddr := args[0]
			argCreatoraddr := args[1]
			argPartneraddr := args[2]
			argHashcode := args[3]
			argChannelid := args[7]

			if err != nil {
				return err
			}
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argNumblock, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			decCoin, err := sdk.ParseDecCoin(args[5])
			if err != nil {
				return err
			}
			coinToCreator, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			decCoin, err = sdk.ParseDecCoin(args[6])
			if err != nil {
				return err
			}
			coinHtlc, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()

			msg := types.NewMsgCommitment(
				clientCtx.GetFromAddress().String(),
				argMultisigAddr,
				argCreatoraddr,
				argPartneraddr,
				argHashcode,
				argNumblock,
				&coinToCreator,
				&coinHtlc,
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
