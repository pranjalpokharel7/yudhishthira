package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/p2p"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

const PORT = ":8080"

func StartServer(wlt *wallet.Wallet, chain *blockchain.BlockChain) {
	// uncomment below line for release mode API
	// gin.SetMode(gin.ReleaseMode)

	go p2p.StartServer("3000", chain)
	router := gin.Default()

	// middlewares
	router.Use(CORSMiddleware())

	// block endpoint
	router.GET("/block/last", GetLastBlockResponse(chain))
	router.GET("/block/last/:n", GetLastNBlocksResponse(chain))

	// item endpoint
	router.GET("/item/history/:itemhash", GetItemTransactionHistoryResponse(chain))
	router.GET("/item/last-block/:itemhash", GetLastBlockWithItemResponse(chain))
	router.GET("/item/calculate-hash/:itemid", CalculateItemHash())
	router.GET("/item/owner/:itemhash", GetItemOwner(chain))

	// general wallet endpoint
	router.GET("/wallet/info/:address", GetWalletInfoResponse(chain))
	router.GET("/wallet/items/:address", GetWalletOwnedItemsResponse(chain))

	// personal wallet endpoint, TODO: combine with generalized wallet above
	router.GET("/my-wallet/address", GetMyWalletAddressResponse(wlt))
	router.GET("/my-wallet/info", GetMyWalletInfoResponse(wlt, chain))
	router.GET("/my-wallet/items", GetMyWalletInfoResponse(wlt, chain))

	// transaction endpoint
	router.GET("/transaction/last/:n", GetLastNTxsResponse(chain))
	router.POST("/transaction/new", PostNewTransaction(wlt, chain))
	router.POST("/transaction/coinbase", PostCoinbaseTransaction(wlt, chain))

	// token verification endpoint
	router.GET("/token/sign/:token", SignToken(wlt))
	router.POST("/token/verify", VerifyToken())

	router.Run(PORT)
}
