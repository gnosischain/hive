package tests

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/hive/hivesim"
	"github.com/golang-jwt/jwt/v4"
)

func TestGenerateNewOnew(t *testing.T) {

}

func TestMain(m *testing.M) {
	suite := hivesim.Suite{
		Name:        "my-suite",
		Description: "This test suite performs some tests.",
	}
	// Setup code here
	setup(suite)
	// hivesim.MustRunSuite(hivesim.New(), suite)

	// Run the tests
	m.Run()

	// Teardown code here, if needed
	teardown(suite)

	// Exit with the code returned by m.Run()
	// os.Exit(code)
}

func setup(suite hivesim.Suite) {
	// Code to run before all tests
	println("Setting up before all tests. Debug mode: %s", os.Getenv("HIVE_DEBUG"))
	if os.Getenv("HIVE_DEBUG") != "false" {
		println("Hive debug mode enabled")
		os.Setenv("HIVE_SIMULATOR", "http://127.0.0.1:3000")
		os.Setenv("HTTP_PROXY", "http://0.0.0.0:8089")
		os.Setenv("BASE_URL", "http://"+"192.168.3.49"+":8545/")
		os.Setenv("ENGINE_URL", "http://"+"192.168.3.49"+":8551/")
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
			Files:       map[string]string{"/genesis.json": "../genesis.json"},
			Run:         runMyClientTest,
		})
		hivesim.MustRunSuite(hivesim.New(), suite)
	}
}

func teardown(suite hivesim.Suite) {
	// Code to run after all tests
	println("Tearing down after all tests")
	// Run the tests. This waits until all tests of the suite
	// have executed.
	// hivesim.MustRunSuite(hivesim.New(), suite)
}

func runMyTest(t *hivesim.T) {
	// write your test code here
}

func runMyClientTest(t *hivesim.T, c *hivesim.Client) {
	s := jwtAuth(hivesim.ENGINEAPI_JWT_SECRET)
	println(s)
	os.Setenv("HIVE_SIMULATOR", "http://127.0.0.1:3000")
	os.Setenv("AUTH_HEADER", "Bearer "+s)
	result := c.RPC().Call("eth_getBlockByNumber", "latest", true)
	generateRawTransaction("0x0202020202020202020202020202020202020202", "ff804d09c833619af673fa99c92ae506d30ff60f37ad41a3d098dcf714db1e4a", "100000", "21000", "0x1", "342770c0", "10203")
	println(result)
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

func TestExample(t *testing.T) {
	t.Log("Running TestExample")
}

func TestTestGenerateNew_the_test(t *testing.T) {
	t.Run("the-test3", func(t *testing.T) {
		runMyTest(&hivesim.T{})
	})
}

func TestTestGenerateNew_the_test_2(t *testing.T) {
	t.Run("the-test-24", func(t *testing.T) {
		runMyClientTest(&hivesim.T{}, &hivesim.Client{})
	})
}

func generateRawTransaction(toAddressString string, privateKeyString string, valueString string, gasLimitString string, nonceString string, gasPriceString string, chainIDString string) {
	// Replace with your own private key
	privateKeyHex := string(privateKeyString)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get the public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to assert public key type")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Public key: %s\n", fromAddress)

	// Get the nonce for the account
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatalf("Failed to get nonce: %v", err)
	// }

	// Set the transaction parameters
	toAddress := common.HexToAddress(toAddressString)        // Replace with recipient address
	value, _ := new(big.Int).SetString(valueString, 10)      // Amount in wei (1 ETH)
	gasLimit, _ := strconv.ParseUint(gasLimitString, 10, 64) // Convert gasLimitString to uint64
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatalf("Failed to suggest gas price: %v", err)
	// }

	nonce, _ := strconv.ParseUint(nonceString, 10, 64)
	gasPrice, _ := new(big.Int).SetString(gasPriceString, 10)
	chainID, _ := new(big.Int).SetString(chainIDString, 10)
	// Create the transaction
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", signedTx.Hash().Hex())
	// Encode the signed transaction to hexadecimal
	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
	}
	rawTxHex := hex.EncodeToString(rawTxBytes)

	// Print the raw transaction hash
	fmt.Printf("Raw transaction hash: %s\n", rawTxHex)
}
