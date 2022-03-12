package api

import (
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

type ErrorJSON struct {
	ErrorMsg string `json:"error"`
}

const PORT = ":8080"

func GetLastBlockResponse(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		lastBlock := chain.GetLastBlock()
		c.JSON(200, lastBlock)
	}
	return fn
}

func GetLastBlockWithItem(chain *blockchain.BlockChain) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		itemHashString := c.Param("itemhash")
		itemHash, err := hex.DecodeString(itemHashString)
		if err != nil {
			c.JSON(400, ErrorJSON{ErrorMsg: "provided hash can not be decoded"})
		}
		lastBlock, _, err := chain.GetLastBlockWithItem(itemHash)
		if err != nil {
			c.JSON(404, ErrorJSON{ErrorMsg: "item with hash not found"})
		}
		c.JSON(200, lastBlock)
	}
	return fn
}

// GET Requests
func GetBlockHeight()                           {}
func GetTransactionByHash(txID string)          {}
func GetTransactionsByBlock(blockHash string)   {}
func GetBlocks(start uint64, end uint64)        {}
func GetAddressFunds(walletAddress string)      {}
func GetItemTransactionHistory(itemHash string) {}

// POST Requests
func PostTransactionRequest(addrFrom string, addrTo string) {}

func StartServer(chain *blockchain.BlockChain) {
	router := gin.Default()

	router.GET("/lastblock", GetLastBlockResponse(chain))
	router.GET("/last-item-block/:itemhash", GetLastBlockWithItem(chain))

	router.Run(PORT)
}
