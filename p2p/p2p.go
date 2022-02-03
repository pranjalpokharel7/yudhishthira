package p2p

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"

	net "github.com/libp2p/go-libp2p-net"

	//internal inports
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

func MakeBasicHost(listenPort int, secio bool, randSeed int64) (host.Host, error) {
	var r io.Reader
	if randSeed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randSeed))
	}
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// build host addr
	hostAddress, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/network/%s", basicHost.ID().Pretty()))
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddress)
	log.Printf("I am %s\n", fullAddr)

	return basicHost, nil
}

func HandleStream(s net.Stream, bChain *blockchain.BlockChain, mutex *sync.Mutex) {
	log.Println("Got new stream")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw, bChain, mutex)
	go writeData(rw, bChain, mutex)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter, bChain *blockchain.BlockChain, mutex *sync.Mutex) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			chain := make([]blockchain.Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(bChain.Blocks) {
				bChain.Blocks = chain
				bytes, err := json.MarshalIndent(chain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter, bChain *blockchain.BlockChain, mutex *sync.Mutex) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			chain := make([]blockchain.Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(bChain.Blocks) {
				bChain.Blocks = chain
				bytes, err := json.MarshalIndent(chain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}
