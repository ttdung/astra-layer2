package cli

import (
	"context"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListFwdcommitment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-fwdcommitment",
		Short: "list all fwdcommitment",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllFwdcommitmentRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.FwdcommitmentAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowFwdcommitment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fwdcommitment [index]",
		Short: "shows a fwdcommitment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			params := &types.QueryGetFwdcommitmentRequest{
				Index: argIndex,
			}

			res, err := queryClient.Fwdcommitment(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
