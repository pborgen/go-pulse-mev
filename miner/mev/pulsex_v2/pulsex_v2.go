package pulsex_v2

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

)

func PopulateReserves(client *ethclient.Client, pairAddress common.Address) (*struct {
	Reserve0 big.Int
	Reserve1 big.Int
}, error) {
	const maxRetries = 5
	var lastErr error

	for i := 0; i < maxRetries; i++ {
	
		contract, err := NewAbiUniswapv2pairCaller(pairAddress, client)

		if err != nil {
			lastErr = err
			continue
		}

		reserves, err := contract.GetReserves(nil)
		if err != nil {
			lastErr = err
			continue
		}

		outStruct := new(struct {
			Reserve0 big.Int
			Reserve1 big.Int
		})

		outStruct.Reserve0 = *reserves.Reserve0
		outStruct.Reserve1 = *reserves.Reserve1

		return outStruct, nil
	}

	return nil, fmt.Errorf("failed after %d retries, last error: %v", maxRetries, lastErr)
}