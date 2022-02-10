package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"

	//internal inports
	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

var bChain *blockchain.BlockChain
var mutex sync.Mutex

const (
	UNNAMED              = 0x0    // not a full node, It may not be able to provide any data except for the transactions it originates.
	NODE_NETWORK         = 0x01   // full node, can be asked for full blocks
	NODE_NETWORK_LIMITED = 0x0400 // same as node network but node has at least last (some amount of blocks to be decided)
)

// a version message that contains the version of the chain in the node
type Version struct {
	version     int32
	Timestamp   uint64
	AddressFrom string
}

func MakeBasicHost(listenPort int, secio bool, randSeed int64) (host.Host, error) {

	var r io.Reader
	if randSeed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randSeed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	if secio {
		opts = append(opts, libp2p.NoSecurity)
	}

	return libp2p.New(opts...)
}

func GetHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ha.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func StartListener(ctx context.Context, ha host.Host, listenPort int, insecure bool) {
	var fullAddr = GetHostAddress(ha)
	log.Printf("Identity: %s", fullAddr)

	// set a stream handler on the host
	ha.SetStreamHandler(
		"echo/1/0.0", func(s network.Stream) {
			log.Print("Listener received new stream")
			err := doEcho(s)
			if err != nil {
				log.Println(err)
				s.Reset()
			} else {
				s.Close()
			}

		})
}

func RunSender(ctx context.Context, ha host.Host, targetPeer string) {
	fullAddr := GetHostAddress(ha)
	log.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("sender received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(targetPeer)
	if err != nil {
		log.Println(err)
		return
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Println(err)
		return
	}

	peerid, err := peer.Decode(pid)
	if err != nil {
		log.Println(err)
		return
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", pid))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, peerstore.PermanentAddrTTL)

	log.Println("sender opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/echo/1.0.0")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("sender saying hello")
	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		log.Println(err)
		return
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("read reply: %q\n", out)
}

// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}
