package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

const PORT = ":8080"

func StartServer(wlt *wallet.Wallet, chain *blockchain.BlockChain) {
	// uncomment below line for release mode API
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// middlewares
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

	// transaction endpoint
	router.GET("/transaction/last/:n", GetLastNTxsResponse(chain))
	router.POST("/transaction/new", PostNewTransaction(wlt, chain))
	router.POST("/transaction/coinbase", PostCoinbaseTransaction(wlt, chain))

	router.Run(PORT)
}