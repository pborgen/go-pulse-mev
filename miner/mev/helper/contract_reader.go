package helper

import (
	"fmt"

	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractReader struct {
    client  *ethclient.Client
    abi     abi.ABI
    address common.Address
}

func NewContractReader(contractAddr common.Address, abiJSON string) (*ContractReader, error) {
    parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
    if err != nil {
        return nil, fmt.Errorf("failed to parse ABI: %v", err)
    }

    myClient, err := GetClient()
    if err != nil {
        return nil, fmt.Errorf("failed to dial client: %v", err)
    }

    return &ContractReader{
        abi:     parsedABI,
        client:  myClient,
        address: contractAddr,
    }, nil
}


