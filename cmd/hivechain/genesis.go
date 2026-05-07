package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/exp/slices"
)

var initialBalance, _ = new(big.Int).SetString("1000000000000000000000000000000000000", 10)

const (
	genesisBaseFee  = params.InitialBaseFee
	blocktimeSec    = 10 // hard-coded in core.GenerateChain
	defaultGasLimit = 100_000_000
)

// Ethereum mainnet forks in order of introduction.
var (
	allForkNames = append(preMergeForkNames, posForkNames...)
	lastFork     = allForkNames[len(allForkNames)-1]

	// these are block-number based:
	preMergeForkNames = []string{
		"homestead",
		"tangerinewhistle",
		"spuriousdragon",
		"byzantium",
		"constantinople",
		"petersburg",
		"istanbul",
		"muirglacier",
		"berlin",
		"london",
		"arrowglacier",
		"grayglacier",
		"merge",
	}

	// forks after the merge are timestamp-based:
	posForkNames = []string{
		"shanghai",
		"cancun",
		"prague",
		"osaka",
	}
)

// Gnosis-specific AuRa addresses
var (
	eip1559FeeCollectorAddress   = common.HexToAddress("0x1559000000000000000000000000000000000000")
	blockRewardContractAddress   = common.HexToAddress("0x2000000000000000000000000000000000000001")
	randomnessContractAddress    = common.HexToAddress("0x3000000000000000000000000000000000000001")
	blockGasLimitContractAddress = common.HexToAddress("0x4000000000000000000000000000000000000001")
	registrarAddress             = common.HexToAddress("0x6000000000000000000000000000000000000000")
	withdrawalContractAddress    = common.HexToAddress("0xbabe2bed00000000000000000000000000000003")
	depositContractAddress       = common.HexToAddress("0xbabe2bed00000000000000000000000000000003")
	blockRewardContractCode      = common.FromHex("0x608060405234801561000f575f5ffd5b5060043610610029575f3560e01c8063f91c28981461002d575b5f5ffd5b61004061003b3660046100c0565b610057565b60405161004e92919061012a565b60405180910390f35b60608060405160408152606060208201525f60408201525f6060820152608081f35b5f5f83601f840112610089575f5ffd5b5081356001600160401b0381111561009f575f5ffd5b6020830191508360208260051b85010111156100b9575f5ffd5b9250929050565b5f5f5f5f604085870312156100d3575f5ffd5b84356001600160401b038111156100e8575f5ffd5b6100f487828801610079565b90955093505060208501356001600160401b03811115610112575f5ffd5b61011e87828801610079565b95989497509550505050565b604080825283519082018190525f9060208501906060840190835b8181101561016c5783516001600160a01b0316835260209384019390920191600101610145565b5050838103602080860191909152855180835291810192508501905f5b818110156101a7578251845260209384019390920191600101610189565b5091969550505050505056fea2646970667358221220e16faf108632709dd9a4a2cec3d61ceeff2919e9ea4162aed99fe5502ef7b8c264736f6c63430008220033")
	withdrawalContractCode       = common.FromHex("0x6080604052348015600e575f5ffd5b50600436106026575f3560e01c806379d0c0bc14602a575b5f5ffd5b603c60353660046082565b5050505050565b005b5f5f83601f840112604d575f5ffd5b5081356001600160401b038111156062575f5ffd5b6020830191508360208260051b8501011115607b575f5ffd5b9250929050565b5f5f5f5f5f606086880312156095575f5ffd5b8535945060208601356001600160401b0381111560b0575f5ffd5b60ba88828901603e565b90955093505060408601356001600160401b0381111560d7575f5ffd5b60e188828901603e565b96999598509396509294939250505056fea2646970667358221220d044cca4aa2b520b54f7d26bdd3e92395a0c5c95d868a809fc97544f0071db4064736f6c63430008220033")
)

// gnosisAuraConfig returns the AuRa consensus config for the Gnosis test chain.
func gnosisAuraConfig() *params.AuRaConfig {
	var (
		stepDuration            = uint64(5)
		blockReward             = uint64(0)
		maxUncleCountTransition = uint64(0)
		maxUncleCount           = uint(0)
		posdaoTransition        = uint64(0)
	)

	dummyValidator := common.HexToAddress("0x7435ed30A8b4AEb0877CEf0c6E8cFFe834eb865f")
	return &params.AuRaConfig{
		StepDuration:                &stepDuration,
		BlockReward:                 &blockReward,
		MaximumUncleCountTransition: &maxUncleCountTransition,
		MaximumUncleCount:           &maxUncleCount,

		Validators: &params.ValidatorSetJson{
			List: []common.Address{dummyValidator},
		},
		RandomnessContractAddress: map[uint64]common.Address{
			0: randomnessContractAddress,
		},
		PosdaoTransition: &posdaoTransition,
		BlockGasLimitContractTransitions: map[uint64]common.Address{
			0: blockGasLimitContractAddress,
		},
		Registrar:                  &registrarAddress,
		Eip1559FeeCollector:        &eip1559FeeCollectorAddress,
		BlockRewardContractAddress: &blockRewardContractAddress,
		WithdrawalContractAddress:  &withdrawalContractAddress,
	}
}

// createChainConfig creates a chain configuration.
func (cfg *generatorConfig) createChainConfig() *params.ChainConfig {
	chaincfg := new(params.ChainConfig)

	chainid, _ := new(big.Int).SetString("3503995874084926", 10)
	chaincfg.ChainID = chainid

	// Set consensus algorithm.
	chaincfg.Aura = gnosisAuraConfig()

	// Set deposit contract address.
	chaincfg.DepositContractAddress = depositContractAddress

	// Apply forks.
	forks := cfg.forkBlocks()
	if _, ok := forks["cancun"]; ok {
		chaincfg.BlobScheduleConfig = new(params.BlobScheduleConfig)
	}
	for fork, b := range forks {
		timestamp := cfg.blockTimestamp(b)

		switch fork {
		// number-based forks
		case "homestead":
			chaincfg.HomesteadBlock = new(big.Int).SetUint64(b)
		case "tangerinewhistle":
			chaincfg.EIP150Block = new(big.Int).SetUint64(b)
		case "spuriousdragon":
			chaincfg.EIP155Block = new(big.Int).SetUint64(b)
			chaincfg.EIP158Block = new(big.Int).SetUint64(b)
		case "byzantium":
			chaincfg.ByzantiumBlock = new(big.Int).SetUint64(b)
		case "constantinople":
			chaincfg.ConstantinopleBlock = new(big.Int).SetUint64(b)
		case "petersburg":
			chaincfg.PetersburgBlock = new(big.Int).SetUint64(b)
		case "istanbul":
			chaincfg.IstanbulBlock = new(big.Int).SetUint64(b)
		case "muirglacier":
			chaincfg.MuirGlacierBlock = new(big.Int).SetUint64(b)
		case "berlin":
			chaincfg.BerlinBlock = new(big.Int).SetUint64(b)
		case "london":
			chaincfg.LondonBlock = new(big.Int).SetUint64(b)
		case "arrowglacier":
			chaincfg.ArrowGlacierBlock = new(big.Int).SetUint64(b)
		case "grayglacier":
			chaincfg.GrayGlacierBlock = new(big.Int).SetUint64(b)
		case "merge":
			chaincfg.MergeNetsplitBlock = new(big.Int).SetUint64(b)
		// time-based forks
		case "shanghai":
			chaincfg.ShanghaiTime = &timestamp
		case "cancun":
			chaincfg.CancunTime = &timestamp
			chaincfg.BlobScheduleConfig.Cancun = params.DefaultCancunBlobConfig
		case "prague":
			chaincfg.PragueTime = &timestamp
			chaincfg.BlobScheduleConfig.Prague = params.DefaultPragueBlobConfig
		case "osaka":
			chaincfg.OsakaTime = &timestamp
			chaincfg.BlobScheduleConfig.Osaka = params.DefaultOsakaBlobConfig
		default:
			panic(fmt.Sprintf("unknown fork name %q", fork))
		}
	}

	// Special case for merged-from-genesis networks.
	// Need to assign TTD here because the genesis block won't be processed by GenerateChain.
	if chaincfg.MergeNetsplitBlock != nil && chaincfg.MergeNetsplitBlock.Sign() == 0 {
		chaincfg.TerminalTotalDifficulty = cfg.genesisDifficulty()
	}

	return chaincfg
}

func (cfg *generatorConfig) genesisDifficulty() *big.Int {
	if cfg.merged {
		return big.NewInt(0)
	}
	return new(big.Int).Set(params.MinimumDifficulty)
}

// createGenesis creates the genesis block and config.
func (cfg *generatorConfig) createGenesis() *core.Genesis {
	var g core.Genesis
	g.Config = cfg.createChainConfig()

	// Block attributes.
	g.Difficulty = cfg.genesisDifficulty()
	g.ExtraData = []byte("hivechain")
	g.GasLimit = cfg.gasLimit
	zero := new(big.Int)
	if g.Config.IsLondon(zero) {
		g.BaseFee = big.NewInt(genesisBaseFee)
	}

	// Initialize allocation.
	// Here we add balance to known accounts and initialize built-in contracts.
	g.Alloc = make(types.GenesisAlloc)
	for _, acc := range knownAccounts {
		g.Alloc[acc.addr] = types.Account{Balance: initialBalance}
	}
	addCancunSystemContracts(g.Alloc)
	addPragueSystemContracts(g.Alloc)
	addSnapTestContract(g.Alloc)
	addModContracts(g.Alloc)
	addGnosisSystemContracts(g.Alloc)

	return &g
}

func addCancunSystemContracts(ga types.GenesisAlloc) {
	ga[params.BeaconRootsAddress] = types.Account{
		Balance: big.NewInt(42),
		Code:    params.BeaconRootsCode,
	}
}

func addPragueSystemContracts(ga types.GenesisAlloc) {
	ga[params.HistoryStorageAddress] = types.Account{Balance: big.NewInt(1), Code: params.HistoryStorageCode}
	ga[params.WithdrawalQueueAddress] = types.Account{Balance: big.NewInt(1), Code: params.WithdrawalQueueCode}
	ga[params.ConsolidationQueueAddress] = types.Account{Balance: big.NewInt(1), Code: params.ConsolidationQueueCode}
}

func addSnapTestContract(ga types.GenesisAlloc) {
	addr := common.HexToAddress("0x8bebc8ba651aee624937e7d897853ac30c95a067")
	h := common.HexToHash
	ga[addr] = types.Account{
		Balance: big.NewInt(1),
		Nonce:   1,
		Storage: map[common.Hash]common.Hash{
			h("0x01"): h("0x01"),
			h("0x02"): h("0x02"),
			h("0x03"): h("0x03"),
		},
	}
}

const (
	emitAddr      = "0x7dcd17433742f4c0ca53122ab541d0ba67fc27df"
	largeLogsAddr = "0x8dcd17433742f4c0ca53122ab541d0ba67fc27ff"
)

// addModContracts adds the contracts used by block modifiers.
func addModContracts(ga types.GenesisAlloc) {
	ga[common.HexToAddress(emitAddr)] = types.Account{
		Code:    emitCode,
		Balance: new(big.Int),
	}
	ga[common.HexToAddress(largeLogsAddr)] = types.Account{
		Code:    modLargeReceiptCode,
		Balance: new(big.Int),
	}
}

// addGnosisSystemContracts adds the system contracts used by Gnosis.
func addGnosisSystemContracts(ga types.GenesisAlloc) {
	ga[blockRewardContractAddress] = types.Account{
		Balance: big.NewInt(1),
		Code:    blockRewardContractCode,
	}
	ga[withdrawalContractAddress] = types.Account{
		Balance: big.NewInt(1),
		Code:    withdrawalContractCode,
	}
}

// forkBlocks computes the block numbers where forks occur. Forks get enabled based on the
// forkInterval. If the total number of requested blocks (chainLength) is lower than
// necessary, the remaining forks activate on the last chain block.
func (cfg *generatorConfig) forkBlocks() map[string]uint64 {
	lastIndex := cfg.lastForkIndex()
	forks := allForkNames[:lastIndex+1]
	forkBlocks := make(map[string]uint64)

	// If merged chain is specified, schedule all pre-merge forks at block zero.
	if cfg.merged {
		for _, fork := range preMergeForkNames {
			if len(forks) == 0 {
				break
			}
			forkBlocks[fork] = 0
			if forks[0] != fork {
				panic("unexpected fork in allForkNames: " + forks[0])
			}
			forks = forks[1:]
		}
	}
	// Schedule remaining forks according to interval.
	for block := 0; block <= cfg.chainLength && len(forks) > 0; {
		fork := forks[0]
		forks = forks[1:]
		forkBlocks[fork] = uint64(block)
		block += cfg.forkInterval
	}
	// If the chain length cannot accommodate the spread of forks with the chosen
	// interval, schedule the remaining forks at the last block.
	for _, f := range forks {
		forkBlocks[f] = uint64(cfg.chainLength)
	}
	return forkBlocks
}

// lastForkIndex returns the index of the latest enabled for in allForkNames.
func (cfg *generatorConfig) lastForkIndex() int {
	if cfg.lastFork == "" || cfg.lastFork == "frontier" {
		return len(allForkNames) - 1
	}
	index := slices.Index(allForkNames, strings.ToLower(cfg.lastFork))
	if index == -1 {
		panic(fmt.Sprintf("unknown lastFork name %q", cfg.lastFork))
	}
	return index
}

func (cfg *generatorConfig) blockTimestamp(num uint64) uint64 {
	return num * blocktimeSec
}
