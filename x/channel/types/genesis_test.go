package types_test

import (
	"testing"

	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				ChannelList: []types.Channel{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				CommitmentList: []types.Commitment{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				FwdcommitmentList: []types.Fwdcommitment{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated channel",
			genState: &types.GenesisState{
				ChannelList: []types.Channel{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated commitment",
			genState: &types.GenesisState{
				CommitmentList: []types.Commitment{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated fwdcommitment",
			genState: &types.GenesisState{
				FwdcommitmentList: []types.Fwdcommitment{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
