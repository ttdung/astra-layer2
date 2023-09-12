package pubkey

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v12/crypto/ethsecp256k1"
	"sort"
	"strings"
	//"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

type PKAccount struct {
	publicKey ethsecp256k1.PubKey //evm
	//publicKey secp256k1.PubKey
}

func NewPKAccount(pubkey string) (*PKAccount, error) {
	//pubkey = strings.Split(strings.Split(pubkey, "{")[1], "}")[0]

	//key, err := hex.DecodeString(pubkey)
	//if err != nil {
	//	return nil, err
	//}

	key := []byte(pubkey)

	return &PKAccount{
		//publicKey: secp256k1.PubKey{
		publicKey: ethsecp256k1.PubKey{
			Key: key,
		},
	}, nil
}

func (pka *PKAccount) PublicKey() cryptoTypes.PubKey {
	return &pka.publicKey
}

func (pka *PKAccount) AccAddress() types.AccAddress {
	pub := pka.PublicKey()
	addr := types.AccAddress(pub.Address())

	return addr
}

func CreateMulSignAccountFromTwoAccount(account1, account2 cryptoTypes.PubKey,
	threshold int) (string, cryptoTypes.PubKey, error) {
	pks := make([]cryptoTypes.PubKey, 2)
	pks[0] = account1
	pks[1] = account2

	fmt.Println("PKs0========.:", pks[0])
	fmt.Println("PKs1========.:", pks[1])
	sort.Slice(pks, func(i, j int) bool {
		return strings.Compare(pks[i].String(), pks[j].String()) < 0
	})

	fmt.Println("PKs0 sorted========.:", pks[0])

	pk := multisig.NewLegacyAminoPubKey(threshold, pks)
	fmt.Println("Multisig pk Address========.:", pk.Address())
	addr := types.AccAddress(pk.Address())
	return addr.String(), pk, nil
}
