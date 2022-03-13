package api

import (
	"encoding/hex"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

type ErrorJSON struct {
	ErrorMsg string `json:"error"`
}

const PORT = ":8080"

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

// POST Requests

func PostNewTransaction(addrFrom string, addrTo string) {

}

func PostCoinbaseTransaction(walletAddress string) {

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func StartServer(chain *blockchain.BlockChain) {
	router := gin.Default()

	router.Use(CORSMiddleware())
	// block endpoint
	router.GET("/block/last", GetLastBlockResponse(chain))
	router.GET("/block/last/:n", GetLastNBlocksResponse(chain))

	// item endpoint
	router.GET("/item/history/:itemhash", GetItemTransactionHistoryResponse(chain))
	router.GET("/item/last-block/:itemhash", GetLastBlockWithItemResponse(chain))

	// wallet endpoint
	router.GET("/wallet/info/:address", GetWalletInfoResponse(chain))
	router.GET("/wallet/items/:address", GetWalletOwnedItemsResponse(chain)) // get items currently owned by the wallet address

	// misc endpoint
	router.POST("/transaction/new")
	router.POST("/transaction/coinbase")

	router.Run(PORT)
}
