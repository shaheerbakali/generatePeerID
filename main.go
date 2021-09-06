package main

import (
	"fmt"
	"math/big"
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"crypto/sha256"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)


var (
	alphabet   = regexp.MustCompile("^[123456789abcdefghijklmnopqrstuvwxyz]+$")
	numWorkers = runtime.NumCPU()
)

// Key is a vector of bits backed by a Go byte slice.
// First byte is most significant.
// First bit (in each byte) is least significant.
type Key1 []byte


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
	for i := 0; i < numWorkers; i++ {
		go func() {
			err := generateKey(part, keyChan)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	for key := range keyChan {

		fmt.Print("-----------------------------------------------------------\n")
		fmt.Print("gonna calculate distance here:\n")


		targetCID := "bafkreicu6zojf7vrlqdn4pmhjqkknvgw7vncn5adv4kxs3o6vbe3ni5zzi"
		fmt.Print("target CID: ",targetCID,"\n")
		fmt.Print("newly generated peer ID: ",key.PrettyID,"\n")
		println("\n")

		targetCIDByte := []byte(targetCID)
		newPeerID := []byte(key.PrettyID)

		targetCIDByteSHA := sha256.Sum256([]byte(targetCIDByte))
		newPeerIDSHA := sha256.Sum256([]byte(newPeerID))


		//fmt.Print("sha256(b1): ",targetCIDByteSHA,"\n")
		//fmt.Print("sha256(b2): ",newPeerIDSHA,"\n")

		myCID := "12D3KooWMCWJH8AfGTJ7UNT9zyWPqdi8EuSJzjzAXfByELmiqPV8"
		myCIDByte := []byte(myCID)
		myCIDSHA := sha256.Sum256([]byte(myCIDByte))
		distance0 := DistInt(targetCIDByteSHA[:],myCIDSHA[:])
		fmt.Print("distance between my peerID and CID: \n")
		fmt.Print(distance0,"\n")
		println("\n")


		distance1 := DistInt(targetCIDByteSHA[:],newPeerIDSHA[:])
		fmt.Print("distance between new peerID and CID: \n")
		fmt.Print(distance1,"\n")
		println("\n")


		// distance between new CID and other closest nodes
		closestNodes := []string {"12D3KooWDgzqBxyfbkXqV9Ni2ZAqAXQSz4J5HAiEk1VXuiyGVAHA",
			"12D3KooWAb7RpmZxYZExzoBxnShrksuAGEsuG13BTHuotwPRFKyB",
			"12D3KooWAbiMqcwudzf7WWNbZUTkvSYQw4Aj6C9tT6NDdSw8LBnQ",
			"12D3KooWFaMLcc6G1ntbjJU9TcL9f5aBTFVZ2KnSuEt9XC86bKeL",
			"12D3KooWHNYxUc8Km6Fq3F9vEgAMttP7jfkuM2EZ26yHJ2xSJD8N",
			"12D3KooWNNhG9Qzopb3wtytrxpZdRikMgNq6hWinVmuaWFjYCjcZ",
			"12D3KooWJsP4v7WbNBLFri8dNMwUPfM8PdqHMa1TAYmaRMRhDJ4D",
			"QmSKVUFAyCddg2wDUdZVCfvqG5YCwwJTWY1HRmorebXcKG",
			"QmY2QELCTEFwcyfSVrXTK4Jmj3KzTptWtpKLpatAPBAyxw",
			"QmXgUGND4x9gmNn5uddq3mfnJFATcPkaHdUAs8iCjDLczB",
			"QmcT1nzgdr8yNgbDx83WZXkxaCs7QpuvzcR3k6S1R2BFPA",
			"QmcoEPpAg5VyB2qTve3BeJ3nc7EUd8Zjq8RjPg8jo1Q7TG",
			"12D3KooWMoivkxPw2Tm35agWAhPf3hFM6u3ZyNeQgE8Dnsbz8NmN"}

		for _, element := range closestNodes{
			elementByte := []byte(element)
			elementByteSHA := sha256.Sum256([]byte(elementByte))
			distance2 := DistInt(targetCIDByteSHA[:],elementByteSHA[:])
			fmt.Print("distance between ",key.PrettyID," and ",element,":\n")
			fmt.Print(distance2,"\n")
			println("\n")
		}




		fmt.Print("-----------------------------------------------------------\n")

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
func (k Key1) NormInt() *big.Int {
	return big.NewInt(0).SetBytes(k)
}

func Xor(x, y Key1) Key1 {
	if len(x) != len(y) {
		fmt.Print("len error!\n")
		return x
	}
	z := make(Key1, len(x))
	for i := range x {
		z[i] = x[i] ^ y[i]
	}
	return z
}

func DistInt(x, y []byte) *big.Int {
	return Xor(x, y).NormInt()
}

// XOR allocates a new byte slice with the computed result of XOR(a, b).
func XOR(a, b []byte) []byte {
	if len(a) != len(b) {
		return a
	}

	c := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}

	return c
}

