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

func CmdAcceptfund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acceptfund [creatoraddr] [channelid] [cointo-creator] [hashcode] [numblock] [multisig-addr]",
		Short: "Broadcast message acceptfund",
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

			var argCointoCreator []*sdk.Coin
			arg2 := strings.Split(args[2], ":")

			argCointoCreator = make([]*sdk.Coin, len(arg2))
			for i, coin := range arg2 {
				decCoin, err := sdk.ParseDecCoin(coin)
				if err != nil {
					return err
				}
				c, _ := sdk.NormalizeDecCoin(decCoin).TruncateDecimal()
				argCointoCreator[i] = &c
			}

			argNumblock, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAcceptfund(
				clientCtx.GetFromAddress().String(),
				argCreatoraddr,
				argChannelid,
				argCointoCreator,
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
