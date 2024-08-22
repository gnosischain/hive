package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/hive/hivesim"
	"github.com/golang-jwt/jwt/v4"

	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/hive/simulators/gnosis/go-bdd/config"
	_ "github.com/ethereum/hive/simulators/gnosis/go-bdd/tests"
)

func main() {
	suite := hivesim.Suite{
		Name:        "my-suite",
		Description: "This test suite performs some tests1.",
	}
	// add a plain test (does not run a client)
	suite.Add(hivesim.TestSpec{
		Name:        "the-test",
		Description: "This is an example test case.",
		Run:         runMyTest,
	})
	// add a client test (starts the client)
	suite.Add(hivesim.ClientTestSpec{
		Name:        "the-test-2",
		Description: "This is an example test case.",
		Files:       map[string]string{"/genesis.json": "genesis.json"},
		Run:         runMyClientTest,
	})
	os.Setenv("HIVE_DEBUG", "false")
	os.Setenv("HTTP_PROXY", "")

	// Run the tests. This waits until all tests of the suite
	// have executed.
	hivesim.MustRunSuite(hivesim.New(), suite)
}

func runMyTest(t *hivesim.T) {
	config.PLACEHOLDERS["BASE_URL"] = "http://"
	//tests.TestGetValueFromPlaceholders(t)
	// write your test code here
}

func runMyClientTest(t *hivesim.T, c *hivesim.Client) {
	println("Running client test")
	// http://172.17.0.6:8545/
	os.Setenv("BASE_URL", "http://"+string(c.IP.String())+":8545/")
	os.Setenv("ENGINE_URL", "http://"+string(c.IP.String())+":8551/")
	println("BASE_URL: " + os.Getenv("BASE_URL"))

	s := jwtAuth(hivesim.ENGINEAPI_JWT_SECRET)
	os.Setenv("HIVE_SIMULATOR", "http://127.0.0.1:3000")
	os.Setenv("AUTH_HEADER", "Bearer "+s)
	// Define the command you want to execute
	cmd := exec.Command("go", "test", "./tests", "-test.v", "-test.run", "^TestCancunFeatures$")

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()

	// Handle any errors that occur
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
	generateRawTransaction(c)
	// Print the output from the command
	fmt.Printf("Output:\n%s\n", string(output))
	println("Finished running client test")
}

var ENGINEAPI_JWT_SECRET = [32]byte{0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x65}

func jwtAuth(secret [32]byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": &jwt.NumericDate{Time: time.Now()},
	})
	s, err := token.SignedString(secret[:])
	if err != nil {
		return ""
	}
	return s
}

func generateRawTransaction(c *hivesim.Client) {
	//Connect to Ethereum node
	//  client, err := rpc.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	//  if err != nil {
	//      log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	//  }

	// Replace with your own private key
	privateKeyHex := string("9c647b8b7c4e7c3490668fb6c11473619db80c93704c70893d3813af4090c39c")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get the public key and address
	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatalf("Failed to assert public key type")
	// }
	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce for the account
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatalf("Failed to get nonce: %v", err)
	// }

	// Set the transaction parameters
	toAddress := common.HexToAddress("0x0202020202020202020202020202020202020202") // Replace with recipient address
	value := big.NewInt(1000000000000000000)                                       // Amount in wei (1 ETH)
	gasLimit := uint64(21000)                                                      // Gas limit
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatalf("Failed to suggest gas price: %v", err)
	// }

	nonce := uint64(0x1)
	gasPrice := big.NewInt(0x3B9ACA00)
	chainID := big.NewInt(10203)
	// Create the transaction
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction
	// chainID, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	log.Fatalf("Failed to get network ID: %v", err)
	// }
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	log.Fatalf("Failed to send transaction: %v", err)
	// }
	fmt.Printf("Transaction sent! Hash: %s\n", signedTx.Hash().Hex())
}
