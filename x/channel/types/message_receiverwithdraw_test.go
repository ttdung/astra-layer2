package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dungtt-astra/astra/v3/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgReceiverwithdraw_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgReceiverwithdraw
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgReceiverwithdraw{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgReceiverwithdraw{
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
