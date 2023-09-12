package keeper

import (
	"github.com/dungtt-astra/astra/v3/x/channel/types"
)

var _ types.QueryServer = Keeper{}
