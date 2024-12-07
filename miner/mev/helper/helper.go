package helper

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)


func GetClient() (*ethclient.Client, error) {

	client, err := ethclient.Dial("https://rpc.pulsechain.com")
        
	if err != nil {
		log.Error("Error in GetClient")
		return nil, err
	} else {
		log.Info("Success! you are connected to the Network")
	}

	return client, nil
}

func GetClientWebSocket() *ethclient.Client {

	client, err := ethclient.Dial("ws://rpc.pulsechain.com")

	if err != nil {
		log.Error("Error in GetClient")
	} else {
		log.Info("Success! you are connected to the Network")
	}

	return client
}