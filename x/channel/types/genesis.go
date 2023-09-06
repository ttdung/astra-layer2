package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ChannelList:       []Channel{},
		CommitmentList:    []Commitment{},
		FwdcommitmentList: []Fwdcommitment{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in channel
	channelIndexMap := make(map[string]struct{})

	for _, elem := range gs.ChannelList {
		index := string(ChannelKey(elem.Index))
		if _, ok := channelIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for channel")
		}
		channelIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in commitment
	commitmentIndexMap := make(map[string]struct{})

	for _, elem := range gs.CommitmentList {
		index := string(CommitmentKey(elem.Index))
		if _, ok := commitmentIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for commitment")
		}
		commitmentIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in fwdcommitment
	fwdcommitmentIndexMap := make(map[string]struct{})

	for _, elem := range gs.FwdcommitmentList {
		index := string(FwdcommitmentKey(elem.Index))
		if _, ok := fwdcommitmentIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for fwdcommitment")
		}
		fwdcommitmentIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
