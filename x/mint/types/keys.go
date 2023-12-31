package types

const (
	ModuleName = "mint"

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	QueryParameters           = "parameters"
	QueryInflation            = "inflation"
	QueryAnnualProvisions     = "annual_provisions"
	QueryTotalMintedProvision = "total_minted_provision"
	QueryBlockProvision       = "block_provision"
	QueryCirculatingSupply    = "circulating_supply"
	QueryBondedRatio          = "bonded_ratio"
)

const (
	prefixMinter = iota
	prefixTotalMintedProvision
)

// MinterKey is the key to use for the keeper store.
var (
	MinterKey                     = []byte{prefixMinter}
	KeyPrefixTotalMintedProvision = []byte{prefixTotalMintedProvision}
)
