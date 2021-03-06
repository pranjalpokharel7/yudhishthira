package api

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/p2p"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

// GET Requests

func GetLastBlockResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		lastBlock := chain.LastBlock()
		c.JSON(200, lastBlock)
	}
	return fn
}

func GetLastBlockWithItemResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		itemHashString := c.Param("itemhash")
		itemHash, err := hex.DecodeString(itemHashString)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "provided hash can not be decoded"})
			return
		}
		lastBlock, _, err := chain.LastBlockWithItem(itemHash)
		if err != nil {
			c.JSON(404, ErrorJSON{ErrorMsg: "item with hash not found"})
			return
		}
		c.JSON(200, lastBlock)
	}
	return fn
}

func GetLastNBlocksResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "invalid height provided: can not be parsed as integer"})
			return
		}
		if n < 0 {
			c.JSON(400, ErrorJSON{ErrorMsg: "negative height provided: block height can only be positive"})
			return
		}
		lastNBlocks := chain.GetLastNBlocks(uint64(n))
		c.JSON(200, lastNBlocks)

	}
	return fn
}

func GetLastNTxsResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "invalid number provided: can not be parsed as integer"})
			return
		}
		if n < 0 {
			c.JSON(400, ErrorJSON{ErrorMsg: "negative number provided: tx count can only be positive"})
			return
		}
		lastNBlocks := chain.GetLastNTxs(uint64(n))
		c.JSON(200, lastNBlocks)
	}
	return fn
}

func GetItemTransactionHistoryResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		itemHashString := c.Param("itemhash")
		itemHash, err := hex.DecodeString(itemHashString)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "provided hash can not be decoded"})
			return
		}
		itemTxHistory := chain.TxsIncludingItem(itemHash)
		c.JSON(200, itemTxHistory)
	}
	return fn
}

func GetWalletInfoResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		walletAddress := c.Param("address")
		coinbaseTxs, err := chain.WalletCoinBaseTxs(walletAddress)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		minedBlocks, err := chain.WalletMinedBlocks(walletAddress)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		walletInfo := map[string]interface{}{
			"coinbase_txs": coinbaseTxs,
			"mined_blocks": minedBlocks,
		}
		c.JSON(200, walletInfo)
	}
	return fn
}

func GetWalletOwnedItemsResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		walletAddress := c.Param("address")
		ownedItems, err := chain.WalletOwnedItems(walletAddress)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		walletInfo := map[string]interface{}{
			"owned_items": ownedItems,
		}
		c.JSON(200, walletInfo)
	}
	return fn
}

func GetMyWalletInfoResponse(wlt *wallet.Wallet, chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		coinbaseTxs, err := chain.WalletCoinBaseTxs(string(wlt.Address))
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		minedBlocks, err := chain.WalletMinedBlocks(string(wlt.Address))
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		walletInfo := map[string]interface{}{
			"coinbase_txs": coinbaseTxs,
			"mined_blocks": minedBlocks,
		}
		c.JSON(200, walletInfo)
	}
	return fn
}

// TODO: make a generalized function for this
func GetMyWalletOwnedItemsResponse(wlt *wallet.Wallet, chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ownedItems, err := chain.WalletOwnedItems(string(wlt.Address))
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		walletInfo := map[string]interface{}{
			"owned_items": ownedItems,
		}
		c.JSON(200, walletInfo)
	}
	return fn
}

func GetMyWalletAddressResponse(wlt *wallet.Wallet) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		walletAddress := string(wlt.Address)
		walletPubkeyHash, err := wallet.PubKeyHashFromAddress(string(wlt.Address))
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		walletPublicKey, err := wallet.PublicKeyToBytes(&wlt.PublicKey)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
			return
		}
		walletPubkeyHashHex := hex.EncodeToString(walletPubkeyHash)
		walletPublicKeyHex := hex.EncodeToString(walletPublicKey)
		walletAddressInfo := map[string]interface{}{
			"address":         walletAddress,
			"public_key":      walletPublicKeyHex,
			"public_key_hash": walletPubkeyHashHex,
		}
		c.JSON(200, walletAddressInfo)
	}
	return fn
}

func PostNewTransaction(wlt *wallet.Wallet, chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		newTxData := NewTxFormInput{}
		if err := c.BindJSON(&newTxData); err != nil {
			c.AbortWithError(400, err)
			return
		}

		itemHash, err := hex.DecodeString(newTxData.ItemHash)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad item hash: could not decode item hex string"})
			return
		}
		newTx, err := blockchain.NewTransaction(wlt, newTxData.Destination, itemHash, newTxData.Amount, chain)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}
		for _, nodeAddress := range p2p.KnownNodes {
			p2p.SendTx(nodeAddress, *newTx)
		}

		c.JSON(200, newTx)
	}
	return fn
}

func PostCoinbaseTransaction(wlt *wallet.Wallet, chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		coinBaseTxData := CoinBaseTxFormInput{}
		if err := c.BindJSON(&coinBaseTxData); err != nil {
			c.AbortWithError(400, err)
			return
		}
		itemHash := sha256.Sum256([]byte(coinBaseTxData.ItemHash))
		coinBaseTx, err := blockchain.CoinBaseTransaction(wlt, itemHash[:], coinBaseTxData.Amount, chain)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}

		for _, nodeAddress := range p2p.KnownNodes {
			p2p.SendTx(nodeAddress, *coinBaseTx)
		}

		c.JSON(200, coinBaseTx)
	}
	return fn
}

func VerifyToken() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		signedTokenData := TokenVerifyModel{}
		if err := c.BindJSON(&signedTokenData); err != nil {
			c.AbortWithError(400, err)
			return
		}

		hashedOriginalToken := sha256.Sum256([]byte(signedTokenData.OriginalToken))
		signedToken, err := hex.DecodeString(signedTokenData.SignedToken)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}

		pubKeyHex := signedTokenData.PublicKey
		pubKeyBytes, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}
		publicKey, err := wallet.BytesToPublicKey(pubKeyBytes)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}

		err = rsa.VerifyPSS(publicKey, crypto.SHA256, hashedOriginalToken[:], signedToken, nil)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}
		verifiedTokenJSON := map[string]interface{}{
			"verified": true,
		}
		c.JSON(200, verifiedTokenJSON)
	}
	return fn
}

func SignToken(wlt *wallet.Wallet) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		tokenData := c.Param("token")
		hashedToken := sha256.Sum256([]byte(tokenData))
		signedToken, err := rsa.SignPSS(rand.Reader, &wlt.PrivateKey, crypto.SHA256, hashedToken[:], nil)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}
		signedTokenJSON := map[string]interface{}{
			"signed_token": hex.EncodeToString(signedToken),
		}
		c.JSON(200, signedTokenJSON)
	}
	return fn
}

func CalculateItemHash() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		itemID := c.Param("itemid")
		itemHash := sha256.Sum256([]byte(itemID))
		itemHashJSON := map[string]interface{}{
			"item_hash": hex.EncodeToString(itemHash[:]),
		}
		c.JSON(200, itemHashJSON)
	}
	return fn
}

func GetItemOwner(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		itemHashString := c.Param("itemhash")
		itemHash, err := hex.DecodeString(itemHashString)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "provided hash can not be decoded"})
			return
		}
		lastTxWithItem, err := blockchain.LastTxWithItem(chain, itemHash)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "provided hash can not be decoded"})
			return
		}
		txOwnerInfo := map[string]interface{}{
			"item_owner": lastTxWithItem.BuyerHash,
		}
		c.JSON(200, txOwnerInfo)
	}
	return fn
}

// TODO: 1. check validity of received transactions, 2. check if transactions exist previously in blockchain
func PostMineBlock(chain *blockchain.BlockChain, wlt *wallet.Wallet) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var txModelPool []TransactionsModel
		if err := c.BindJSON(&txModelPool); err != nil {
			c.AbortWithError(400, err)
			return
		}

		var txPool []blockchain.Tx
		for _, txModel := range txModelPool {
			tx, err := ModelToTx(txModel)
			if err != nil {
				c.AbortWithError(400, err)
				return
			}
			txPool = append(txPool, *tx)
		}
		newBlock := blockchain.CreateBlock()
		newBlock.AddTransactionsToBlock(txPool)
		newBlock.MineBlock(chain, wlt)
		chain.AddBlock(newBlock)

		// TODO: we clear the memory pool here but edit in later commit to remove only selected transactions
		p2p.MemoryPool = map[string]blockchain.Tx{}

		c.JSON(200, newBlock)
	}
	return fn
}

func GetTxPool() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var txsInPools []blockchain.Tx
		for _, tx := range p2p.MemoryPool {
			txsInPools = append(txsInPools, tx)
		}
		c.JSON(200, txsInPools)
	}
	return fn
}
