package keeper_test

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	testkeeper "github.com/AstraProtocol/astra/channel/testutil/keeper"
	"github.com/AstraProtocol/astra/channel/x/channel/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ChannelKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)
	hash := sha256.Sum256([]byte("bcdf"))
	fmt.Println("hashcode:", base64.StdEncoding.EncodeToString(hash[:]))

	require.EqualValues(t, params, k.GetParams(ctx))
}
