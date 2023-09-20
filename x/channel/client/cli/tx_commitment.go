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

			var coinToCreator, coinHtlc []*sdk.Coin

			arg5 := strings.Split(args[5], ":")
			coinToCreator = make([]*sdk.Coin, len(arg5))
			for i, coin := range arg5 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				coinToCreator[i] = &c
			}

			arg6 := strings.Split(args[6], ":")
			coinHtlc = make([]*sdk.Coin, len(arg6))
			for i, coin := range arg6 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				coinHtlc[i] = &c
			}

			msg := types.NewMsgCommitment(
				clientCtx.GetFromAddress().String(),
				argMultisigAddr,
				argCreatoraddr,
				argPartneraddr,
				argHashcode,
				argNumblock,
				coinToCreator,
				coinHtlc,
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
