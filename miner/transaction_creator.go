package miner

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

type CustomTxManager struct {
    signer     types.Signer
    privateKey *ecdsa.PrivateKey
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
        mu:         sync.Mutex{},
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

func (ctm *CustomTxManager) CreateTransactions(w *worker, env *environment) map[common.Address][]*txpool.LazyTransaction {
    
    myTransactions := make(map[common.Address][]*txpool.LazyTransaction)
    ctm.mu.Lock()
    defer ctm.mu.Unlock()
    txPool := w.eth.TxPool()
    legacyPool := txPool.GetLegacyPool()

    



    
	nonce, err := ctm.queryNonce(common.HexToAddress(ctm.GetPublicKeyAsString()))
	if err != nil {
		log.Error("Failed to query nonce: %v", err)
	}

    tx := types.NewTransaction(
        nonce,
        common.HexToAddress(ctm.GetPublicKeyAsString()),
        big.NewInt(1000000000000000000), // 1 PLS
        21000,                           // Gas limit
        big.NewInt(20000000000),         // Gas price (20 Gwei)
        nil,                             // Data
    )

    signedTx, err := types.SignTx(tx, ctm.signer, ctm.privateKey)
    if err != nil {
        log.Error("Failed to sign custom transaction", "err", err)
        return nil
    }

    lazyTransaction := &txpool.LazyTransaction{
        Pool:      legacyPool,
        Hash:      tx.Hash(),
        Tx:        signedTx,
        Time:      tx.Time(),
        GasFeeCap: uint256.MustFromBig(tx.GasFeeCap()),
        GasTipCap: uint256.MustFromBig(tx.GasTipCap()),
        Gas:       tx.Gas(),
        BlobGas:   tx.BlobGas(),
        OrderNumber: 1,
    }

    myTransactions[common.HexToAddress(ctm.GetPublicKeyAsString())] = 
        append(myTransactions[common.HexToAddress(ctm.GetPublicKeyAsString())], lazyTransaction)

    return myTransactions
}


func findMev(w *worker, env *environment) {

    // Original transaction processing code...
	filter := txpool.PendingFilter{
		MinTip: w.tip,
	}

	if env.header.BaseFee != nil {
		filter.BaseFee = uint256.MustFromBig(env.header.BaseFee)
	}

	if env.header.ExcessBlobGas != nil {
		filter.BlobFee = uint256.MustFromBig(eip4844.CalcBlobFee(*env.header.ExcessBlobGas))
	}

	filter.OnlyPlainTxs, filter.OnlyBlobTxs = true, false
	
	pendingPlainTxs := w.eth.TxPool().Pending(filter)

    for addr, txs := range pendingPlainTxs {
        // Iterate through each transaction for this address
        for _, tx := range txs {
            // Access transaction details using tx.Tx
            transaction := tx.Tx
            transaction.Data()
            // Example: Log transaction details
            log.Info("Found pending transaction", 
                "from", addr,
                "to", transaction.To(),
                "value", transaction.Value(),
                "gas", transaction.Gas(),
                "gasPrice", transaction.GasPrice())
        }
    }
}

func (ctm *CustomTxManager) queryNonce(address common.Address) (uint64, error) {
    
    client, err := ethclient.Dial("http://localhost:8545") // Replace with your Ethereum node URL
    if err != nil {
        log.Error("Failed to connect to Ethereum client", "err", err)
    }
    
    context := context.Background()

    // Get the latest confirmed nonce
    confirmedNonce, err := client.NonceAt(context, address, nil) // nil means latest block
    if err != nil {
        return 0, fmt.Errorf("failed to get confirmed nonce: %v", err)
    }

    // Get the pending nonce (includes unconfirmed transactions)
    pendingNonce, err := client.PendingNonceAt(context, address)
    if err != nil {
        return 0, fmt.Errorf("failed to get pending nonce: %v", err)
    }

    // Use the higher nonce to ensure we don't reuse any nonce
    if pendingNonce > confirmedNonce {
        return pendingNonce, nil
    }
    return confirmedNonce, nil
}

func removeTransaction(pool *txpool.TxPool, hash common.Hash) error {
    // Check if transaction exists
    if !pool.Has(hash) {
        return fmt.Errorf("transaction %s not found in pool", hash.Hex())
    }

   


    return nil
}
