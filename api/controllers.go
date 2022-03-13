package api

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

type ErrorJSON struct {
	ErrorMsg string `json:"error"`
}

// TODO: add binding validation
type NewTxFormInput struct {
	Destination string `json:"destination" binding:"required"`
	ItemHash    string `json:"item_hash" binding:"required"`
	Amount      uint64 `json:"amount" binding:"required"`
}

type CoinBaseTxFormInput struct {
	ItemHash string `json:"item_hash" binding:"required"`
	Amount   uint64 `json:"amount" binding:"required"`
}

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
		}
		lastBlock, _, err := chain.LastBlockWithItem(itemHash)
		if err != nil {
			c.JSON(404, ErrorJSON{ErrorMsg: "item with hash not found"})
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
		}
		if n < 0 {
			c.JSON(400, ErrorJSON{ErrorMsg: "negative height provided: block height can only be positive"})
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
		}
		if n < 0 {
			c.JSON(400, ErrorJSON{ErrorMsg: "negative number provided: tx count can only be positive"})
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
		}
		minedBlocks, err := chain.WalletMinedBlocks(walletAddress)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad address: could not derive public key hash from address"})
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
		}
		walletInfo := map[string]interface{}{
			"owned_items": ownedItems,
		}
		c.JSON(200, walletInfo)
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
		itemHash, err := hex.DecodeString(coinBaseTxData.ItemHash)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "bad item hash: could not decode item hex string"})
			return
		}
		coinBaseTx, err := blockchain.CoinBaseTransaction(wlt, itemHash, coinBaseTxData.Amount, chain)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: fmt.Sprintf("%v", err)})
			return
		}
		c.JSON(200, coinBaseTx)
	}
	return fn
}
