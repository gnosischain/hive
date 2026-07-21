package main

import (
	"bytes"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestGenerate(t *testing.T) {
	outdir := t.TempDir()
	cfg := generatorConfig{
		txInterval:   1,
		txCount:      10,
		forkInterval: 2,
		chainLength:  30,
		outputDir:    outdir,
		outputs:      outputFunctionNames(),
	}
	cfg, err := cfg.withDefaults()
	if err != nil {
		t.Fatal(err)
	}
	g := newGenerator(cfg)
	if err := g.run(); err != nil {
		t.Fatal(err)
	}

	names, _ := filepath.Glob(filepath.Join(outdir, "*"))
	t.Log("output files:", names)
}

func TestGnosisExecutionAPIConfig(t *testing.T) {
	cfg := generatorConfig{
		merged:       true,
		forkInterval: 3,
		chainLength:  45,
		lastFork:     "prague",
	}
	cfg, err := cfg.withDefaults()
	if err != nil {
		t.Fatal(err)
	}
	genesis := cfg.createGenesis()
	chaincfg := genesis.Config

	if chaincfg.ChainID.Cmp(big.NewInt(100)) != 0 {
		t.Fatalf("chain ID is %v, want 100", chaincfg.ChainID)
	}
	if chaincfg.Aura == nil {
		t.Fatal("AuRa configuration is missing")
	}
	if chaincfg.Ethash != nil || chaincfg.Clique != nil {
		t.Fatal("unexpected non-AuRa consensus configuration")
	}
	if chaincfg.MergeNetsplitBlock == nil || chaincfg.MergeNetsplitBlock.Sign() != 0 {
		t.Fatalf("merge block is %v, want 0", chaincfg.MergeNetsplitBlock)
	}
	if chaincfg.TerminalTotalDifficulty == nil || chaincfg.TerminalTotalDifficulty.Sign() != 0 {
		t.Fatalf("terminal total difficulty is %v, want 0", chaincfg.TerminalTotalDifficulty)
	}
	if chaincfg.DepositContractAddress != depositContractAddress {
		t.Fatalf("deposit contract is %s, want %s", chaincfg.DepositContractAddress, depositContractAddress)
	}
	if chaincfg.GetMinBlobGasPrice() != params.GnosisBlobTxMinBlobGasprice {
		t.Fatalf("minimum blob gas price is %d, want %d", chaincfg.GetMinBlobGasPrice(), params.GnosisBlobTxMinBlobGasprice)
	}
	if chaincfg.GetMaxBlobsPerTransaction() != params.GnosisBlobTxMaxBlobs {
		t.Fatalf("maximum blobs per transaction is %d, want %d", chaincfg.GetMaxBlobsPerTransaction(), params.GnosisBlobTxMaxBlobs)
	}
	if chaincfg.ShanghaiTime == nil || *chaincfg.ShanghaiTime != 0 {
		t.Fatalf("Shanghai timestamp is %v, want 0", chaincfg.ShanghaiTime)
	}
	if chaincfg.CancunTime == nil || *chaincfg.CancunTime != 30 {
		t.Fatalf("Cancun timestamp is %v, want 30", chaincfg.CancunTime)
	}
	if chaincfg.PragueTime == nil || *chaincfg.PragueTime != 60 {
		t.Fatalf("Prague timestamp is %v, want 60", chaincfg.PragueTime)
	}
	if chaincfg.BlobScheduleConfig == nil ||
		chaincfg.BlobScheduleConfig.Cancun != gnosisBlobConfig ||
		chaincfg.BlobScheduleConfig.Prague != gnosisBlobConfig {
		t.Fatal("Gnosis blob schedule is incomplete")
	}

	// A merged-from-genesis block uses the standard PoS nonce/mixHash header
	// layout even though AuRa remains configured for Gnosis system contracts.
	header := genesis.ToBlock().Header()
	if len(header.Signature) != 0 {
		t.Fatalf("merged genesis has %d-byte AuRa signature, want standard PoS header", len(header.Signature))
	}
	if header.Difficulty.Sign() != 0 {
		t.Fatalf("merged genesis difficulty is %v, want 0", header.Difficulty)
	}

	reward := genesis.Alloc[blockRewardContractAddress]
	if !bytes.Equal(reward.Code, blockRewardContractCode) {
		t.Fatal("block reward system contract has unexpected code")
	}
	if _, ok := genesis.Alloc[withdrawalContractAddress]; !ok {
		t.Fatal("withdrawal system contract is missing")
	}
}

func TestGnosisForkEnv(t *testing.T) {
	cfg := generatorConfig{
		merged:       true,
		forkInterval: 3,
		chainLength:  45,
		lastFork:     "prague",
		outputDir:    t.TempDir(),
	}
	cfg, err := cfg.withDefaults()
	if err != nil {
		t.Fatal(err)
	}
	g := newGenerator(cfg)
	if err := g.writeForkEnv(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(cfg.outputDir, "forkenv.json"))
	if err != nil {
		t.Fatal(err)
	}
	var env map[string]string
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatal(err)
	}
	if got := env["HIVE_DEPOSIT_CONTRACT_ADDRESS"]; got != depositContractAddress.Hex() {
		t.Fatalf("deposit contract env is %q, want %q", got, depositContractAddress.Hex())
	}
	if got := env["HIVE_CANCUN_BLOB_BASE_FEE_UPDATE_FRACTION"]; got != "1112826" {
		t.Fatalf("Cancun blob update fraction is %q, want 1112826", got)
	}
}

func TestGenerateGnosisExecutionAPIChain(t *testing.T) {
	outdir := t.TempDir()
	cfg := generatorConfig{
		merged:       true,
		forkInterval: 3,
		chainLength:  9,
		lastFork:     "prague",
		txInterval:   1,
		txCount:      4,
		outputDir:    outdir,
		outputs: []string{
			"genesis",
			"chain",
			"forkenv",
			"headstate",
			"txinfo",
			"accounts",
			"headfcu",
		},
	}
	cfg, err := cfg.withDefaults()
	if err != nil {
		t.Fatal(err)
	}
	if err := newGenerator(cfg).run(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(outdir, "headfcu.json"))
	if err != nil {
		t.Fatal(err)
	}
	var request rpcRequest
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatal(err)
	}
	if request.Method != "engine_forkchoiceUpdatedV3" {
		t.Fatalf("head forkchoice method is %q, want engine_forkchoiceUpdatedV3", request.Method)
	}
	info, err := os.Stat(filepath.Join(outdir, "chain.rlp"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("generated chain.rlp is empty")
	}
}
