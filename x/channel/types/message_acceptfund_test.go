package types

import (
	"testing"

	"github.com/AstraProtocol/astra/channel/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgAcceptfund_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAcceptfund
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAcceptfund{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgAcceptfund{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
