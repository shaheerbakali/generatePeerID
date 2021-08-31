package main

import (
	"fmt"
	//"math/bits"
	//"sync"
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	//kbucket "github.com/libp2p/go-libp2p-kbucket"
)


var (
	alphabet   = regexp.MustCompile("^[123456789abcdefghijklmnopqrstuvwxyz]+$")
	numWorkers = runtime.NumCPU()
)

// Key stores PrettyID containing desired substring at Index
type Key struct {
	PrettyID string
	Index    int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf(`
This tool generates IPFS public and private keypair until it finds public key
which contains required substring. Keys are stored in local directory. If you
like one of them, you can move it to ~/.ipfs/keystore/ to use it with IPFS.
Usage:
	%s {part}
		For fast results suggested length of public key part is 4-5 characters
`, os.Args[0])
		os.Exit(1)
	}
	part := strings.ToLower(os.Args[1])
	if !alphabet.MatchString(part) {
		fmt.Println("{part} must match the alphabet:", alphabet.String())
		os.Exit(2)
	}
	runtime.GOMAXPROCS(numWorkers)
	keyChan := make(chan Key)
	fmt.Print("keyChan: ",keyChan)
	println("\n")
	for i := 0; i < numWorkers; i++ {
		go func() {
			err := generateKey(part, keyChan)
			fmt.Print("err: ",err)
			println("\n")
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	for key := range keyChan {
		fmt.Print("key: ",key)
		println("\n")
		fmt.Printf(
			"%s\u001b[32m%s\u001b[0m%s\n",
			key.PrettyID[:key.Index],
			key.PrettyID[key.Index:len(part)+key.Index],
			key.PrettyID[len(part)+key.Index:])

	}
}

func generateKey(part string, keyChan chan Key) error {
	for {
		privateKey, publicKey, err := crypto.GenerateEd25519Key(rand.Reader)
		//fmt.Print("privateKey: ",privateKey,"\n")
		//fmt.Print("publicKey: ",publicKey,"\n")
		//println("\n")
		if err != nil {
			return err
		}
		peerID, err := peer.IDFromPublicKey(publicKey)
		if err != nil {
			return err
		}
		prettyID := peerID.Pretty()
		lowerID := strings.ToLower(prettyID)
		idx := strings.Index(lowerID, part)
		if idx == -1 {
			continue
		}
		privateKeyBytes, err := privateKey.Bytes()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(prettyID, privateKeyBytes, 0600)
		if err != nil {
			return err
		}
		keyChan <- Key{
			PrettyID: prettyID,
			Index:    idx,
		}
	}
}

