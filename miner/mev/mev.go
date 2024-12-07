package mev

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
)

type PossibleSandwich struct {
    DexName  string
	DexRouterAddress string
	DexFactoryAddress string

	Token0Address string
	Token1Address string
	PairContractAddress string
}

var possibleSandwich map[string]PossibleSandwich

func init() {
	if possibleSandwich == nil {
		// Read json file to memory
		if data, err := os.ReadFile("possible_sandwich.json"); err != nil {
			panic(err)
		} else {
			json.Unmarshal(data, &possibleSandwich)
		}
	}

}

func IsSandwich(tx *types.Transaction) bool {

	toAddress := tx.To()
	toAddressString := toAddress.String()

	if _, exists := possibleSandwich[toAddressString]; exists {
		return true
	} else {	
		return false
	}
}

func CanSandwich(tx *types.Transaction) (bool, *types.Transaction) {
	return false, nil
}
