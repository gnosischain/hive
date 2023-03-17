package gnosis

import (
	"fmt"

	"github.com/ethereum/hive/simulators/eth2/common/clients"
	"github.com/ethereum/hive/simulators/eth2/common/testnet"
)

func GetWithdrawalsTestnetConfig(tc *testnet.Config, allNodeDefinitions clients.NodeDefinition) *testnet.Config {
	return tc
}

func GetBuilderWithdrawalsTestnetConfig(tc *testnet.Config, allNodeDefinitions clients.NodeDefinitions) *testnet.Config {
	// overridable values:
	//AltairForkEpoch                 *big.Int
	//BellatrixForkEpoch              *big.Int
	//CapellaForkEpoch                *big.Int
	//ValidatorCount                  *big.Int
	//KeyTranches                     *big.Int
	//SlotTime                        *big.Int
	//TerminalTotalDifficulty         *big.Int
	//SafeSlotsToImportOptimistically *big.Int
	//ExtraShares                     *big.Int
	//
	//// Node configurations to launch. Each node as a proportional share of
	//// validators.
	//NodeDefinitions clients.NodeDefinitions
	//Eth1Consensus   execution_config.ExecutionConsensus
	//
	//// Execution Layer specific config
	//InitialBaseFeePerGas     *big.Int
	//GenesisExecutionAccounts map[common.Address]core.GenesisAccount
	//
	//// Builders
	//EnableBuilders bool
	//BuilderOptions []mock_builder.Option
	for i, node := range allNodeDefinitions {
		fmt.Printf("node %d\n", i)
		fmt.Println(node.Chain)
	}

	//tc.Eth1Consensus.HiveParams().Set("gnosis", "true")

	return tc
}
