package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dungtt-astra/astra/v3/x/feeburn/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryTotalFeeBurn() {
	suite.SetupTest()
	ctx := sdk.WrapSDKContext(suite.ctx)
	totalFeeBurn := sdk.NewDec(100000000000)
	suite.app.FeeBurnKeeper.SetTotalFeeBurn(suite.ctx, totalFeeBurn)

	res, err := suite.queryClient.TotalFeeBurn(ctx, &types.QueryTotalFeeBurnRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.TotalFeeBurn.Amount, totalFeeBurn)
}
