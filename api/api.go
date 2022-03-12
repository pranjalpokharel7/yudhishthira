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

// POST Requests

func PostNewTransaction(addrFrom string, addrTo string) {}

func PostCoinbaseTransaction() {}

func StartServer(chain *blockchain.BlockChain) {
	router := gin.Default()

	router.GET("/last-block", GetLastBlockResponse(chain))
	router.GET("/item-history/:itemhash", GetItemTransactionHistoryResponse(chain))
	router.GET("/last-n-blocks/:n", GetLastNBlocksResponse(chain))
	router.GET("/last-item-block/:itemhash", GetLastBlockWithItemResponse(chain))
	router.GET("/wallet/:address", GetWalletInfoResponse(chain))

	router.Run(PORT)
}
