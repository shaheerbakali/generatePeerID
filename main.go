package main

import (
	"fmt"
	"math/big"
	"crypto/rand"
	"io/ioutil"
	"log"
	"runtime"
	"crypto/sha256"
	"os"
	"strconv"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	numWorkers = runtime.NumCPU()
)

// Key1 is a vector of bits backed by a Go byte slice.
type Key1 []byte


// Key stores PrettyID containing desired substring at Index
type Key struct {
	PrettyID string
	Index    int
}

func main() {

	fmt.Printf(`
This tool generates a number of IPFS public and private keypairs randomly. The number is determined by the user input.
Keys are stored in local directory. 
If you like one of them, you can move it to ~/.ipfs/keystore/ to use it with IPFS.
Usage:
	go run .\main.go No. (in terminal)
`,)

	if len(os.Args) < 2 {
		// wrong
		fmt.Print("\n! no input number found !\n")
		os.Exit(1)
	}


	runtime.GOMAXPROCS(numWorkers)
	keyChan := make(chan Key)
	for i := 0; i < numWorkers; i++ {
		go func() {
			err := generateKey(keyChan)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	var distance2 *big.Int
	var NEWcloserPeers []string


	number, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// wrong
		fmt.Print(err,"\n")
	}

	/*
	for i := 0;i< number; i++ {
		fmt.Print("i: ",i,"\n")
	}
	 */


	for key := range keyChan {
		//fmt.Print("\n",key.PrettyID,"\n")
		//println("\n")

		fmt.Print("-----------------------------------------------------------\n")
		fmt.Print("newly generated peer ID: ",key.PrettyID,"\n")
		println("\n")
		//fmt.Print("calculate all distances here:\n")


		targetCID := "bafkreicu6zojf7vrlqdn4pmhjqkknvgw7vncn5adv4kxs3o6vbe3ni5zzi"
		fmt.Print("target CID: ",targetCID,"\n")

		// convert target CID and the newly generated peer ID to bytes
		targetCIDByte := []byte(targetCID)
		newPeerID := []byte(key.PrettyID)

		// hash the target CID and the newly generated peer ID using SHA256
		targetCIDByteSHA := sha256.Sum256([]byte(targetCIDByte))
		newPeerIDSHA := sha256.Sum256([]byte(newPeerID))



		/*
		// to calculate distance between my peer ID and the target CID
		myCID := "12D3KooWMCWJH8AfGTJ7UNT9zyWPqdi8EuSJzjzAXfByELmiqPV8"
		myCIDByte := []byte(myCID)
		myCIDSHA := sha256.Sum256([]byte(myCIDByte))
		distance0 := DistInt(targetCIDByteSHA[:],myCIDSHA[:])
		fmt.Print("distance between my peerID and CID: \n")
		fmt.Print(distance0,"\n")
		println("\n")
		 */




		// to calculate distance between newly generated peer ID and the target CID
		distance1 := DistInt(targetCIDByteSHA[:],newPeerIDSHA[:])
		fmt.Print("distance between newly generated peerID and CID: \n")
		fmt.Print(distance1,"\n")
		println("\n")



		// list of closest nodes of target CID
		closestNodes := []string {
			"12D3KooW9tFbbFUpFiPMrZfX36uTSpJcSZcrVLMqv8YM15fBbMSM",
			"12D3KooWA6eDqwPb2h9BV3aBxaonaW7D7AZvCUcRxdG9YgFXe2B5",
			"12D3KooWAbiMqcwudzf7WWNbZUTkvSYQw4Aj6C9tT6NDdSw8LBnQ",
			"12D3KooWCbjjNF3s8Zx2vd2vcnSViRQi4zKzLadYBFxZ7vY7X1T8",
			"12D3KooWHNYxUc8Km6Fq3F9vEgAMttP7jfkuM2EZ26yHJ2xSJD8N",
			"12D3KooWNT19ZdgJESkXzjEjZXQaXABQ9NeA8imNSrm8ES1phqCF",
			"12D3KooWNNhG9Qzopb3wtytrxpZdRikMgNq6hWinVmuaWFjYCjcZ",
			"12D3KooWJsP4v7WbNBLFri8dNMwUPfM8PdqHMa1TAYmaRMRhDJ4D",
			"12D3KooWLRdiuW625YHLVBAoQwXf3f2UXuAkfsVfwyLv4AyyF8TN",
			"12D3KooWLjd5sJ32CKGefNBRg5HQEgXitGYdtE4KLSWkoyy8HW8q",
			"12D3KooWSK8dtBh7sx5BgGBmaFWshB3sr5vo1vcoi3aP927Nq6Cw",
			"12D3KooWPbh6Qh9yk5HqneY4fpa4ydJN48dT7c4MtDhCv4k9MGhq",
			"QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
			"QmSKVUFAyCddg2wDUdZVCfvqG5YCwwJTWY1HRmorebXcKG",
			"QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
			"QmXgUGND4x9gmNn5uddq3mfnJFATcPkaHdUAs8iCjDLczB",
			"QmcT1nzgdr8yNgbDx83WZXkxaCs7QpuvzcR3k6S1R2BFPA",
			"QmcoEPpAg5VyB2qTve3BeJ3nc7EUd8Zjq8RjPg8jo1Q7TG",
			"QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
			"12D3KooWMoivkxPw2Tm35agWAhPf3hFM6u3ZyNeQgE8Dnsbz8NmN"}
		// to calculate distance between target CID and its closest nodes
		for _, element := range closestNodes{
			// convert CIDs to byte
			elementByte := []byte(element)
			// hash the CIDs
			elementByteSHA := sha256.Sum256([]byte(elementByte))
			// calculate distance between CIDs and target CID
			distance2 = DistInt(targetCIDByteSHA[:],elementByteSHA[:])
		}
		// if distance between newly generated peer ID and the targetCID < distance between targetCID & Closest Nodes
		comparing := distance1.Cmp(distance2)
		if comparing == -1 {
			// so distance1 is shorter
			// means newly generated peerID is much closer
			// save it in a list
			NEWcloserPeers = append(NEWcloserPeers,key.PrettyID)
			fmt.Print("new closer nodes: \n")
			fmt.Print(NEWcloserPeers)
			println("\n")
			if len(NEWcloserPeers) == number {
				os.Exit(1)
			}
		}
		fmt.Print("-----------------------------------------------------------\n")
	}

}

func generateKey(keyChan chan Key) error {
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

