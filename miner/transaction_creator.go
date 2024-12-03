package miner

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type CustomTxManager struct {
	client     *ethclient.Client
    signer     types.Signer
    privateKey *ecdsa.PrivateKey
    nonce      uint64
    mu         sync.Mutex
}

func NewCustomTxManager(chainID *big.Int) *CustomTxManager {
	privateKeyHex := "6af0afaa4552f1b25ae2790be42befd183be8d53271e10cad95a579b27813b44"

	// Convert the hex string to an ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Error("Failed to convert hex to ECDSA: %v", err)
	}

    customTxManager := CustomTxManager{
        signer:     types.NewEIP155Signer(chainID),
        privateKey: privateKey,
        nonce:      0,
    }

	return &customTxManager
}

func (ctm *CustomTxManager) GetPublicKey() *ecdsa.PublicKey {
	return ctm.privateKey.Public().(*ecdsa.PublicKey)
}

// If you need the Ethereum address format instead
func (ctm *CustomTxManager) GetPublicKeyAsString() string {
	publicKey := ctm.privateKey.Public().(*ecdsa.PublicKey)

    return publicKey.X.String()
}

func (ctm *CustomTxManager) CreateTransaction() *types.Transaction {
    ctm.mu.Lock()
    defer ctm.mu.Unlock()

	nonce, err := ctm.queryNonce(common.HexToAddress(ctm.GetPublicKeyAsString()))
	if err != nil {
		log.Error("Failed to query nonce: %v", err)
	}
	ctm.nonce = nonce

    tx := types.NewTransaction(
        ctm.nonce,
        common.HexToAddress(ctm.GetPublicKeyAsString()),
        big.NewInt(1000000000000000000), // 1 ETH
        21000,                           // Gas limit
        big.NewInt(20000000000),         // Gas price (20 Gwei)
        nil,                             // Data
    )

    signedTx, err := types.SignTx(tx, ctm.signer, ctm.privateKey)
    if err != nil {
        log.Error("Failed to sign custom transaction", "err", err)
        return nil
    }

    ctm.nonce++
    return signedTx
}

func (ctm *CustomTxManager) queryNonce(address common.Address) (uint64, error) {
    ctx := context.Background()

    // Get the latest confirmed nonce
    confirmedNonce, err := ctm.client.NonceAt(ctx, address, nil) // nil means latest block
    if err != nil {
        return 0, fmt.Errorf("failed to get confirmed nonce: %v", err)
    }

    // Get the pending nonce (includes unconfirmed transactions)
    pendingNonce, err := ctm.client.PendingNonceAt(ctx, address)
    if err != nil {
        return 0, fmt.Errorf("failed to get pending nonce: %v", err)
    }

    // Use the higher nonce to ensure we don't reuse any nonce
    if pendingNonce > confirmedNonce {
        return pendingNonce, nil
    }
    return confirmedNonce, nil
}