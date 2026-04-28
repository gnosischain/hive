package config

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/hive/simulators/ethereum/engine/config/cancun"
)

func (f *ForkConfig) ConfigGenesis(genesis *core.Genesis) error {
	genesis.Config.TerminalTotalDifficulty = genesis.Difficulty
	genesis.Config.MergeNetsplitBlock = big.NewInt(0)
	if f.ShanghaiTimestamp != nil {
		shanghaiTime := f.ShanghaiTimestamp.Uint64()
		genesis.Config.ShanghaiTime = &shanghaiTime
	}
	if f.CancunTimestamp != nil {
		if err := cancun.ConfigGenesis(genesis, f.CancunTimestamp.Uint64()); err != nil {
			return fmt.Errorf("failed to configure cancun fork: %v", err)
		}
	}

	blockRewardAddr := common.HexToAddress("0x2000000000000000000000000000000000000001")
	randomnessAddr := common.HexToAddress("0x3000000000000000000000000000000000000001")
	blockGasLimitAddr := common.HexToAddress("0x4000000000000000000000000000000000000001")
	registrarAddr := common.HexToAddress("0x6000000000000000000000000000000000000000")
	withdrawalAddr := common.HexToAddress("0xbabe2bed00000000000000000000000000000003")
	feeCollectorAddr := common.HexToAddress("0x1559000000000000000000000000000000000000")
	validatorAddr := common.HexToAddress("0x14747a698Ec1227e6753026C08B29b4d5D3bC484")

	stepDuration := uint64(5)
	blockReward := uint64(0)
	blockRewardContractTransition := uint64(0)
	maximumUncleCountTransition := uint64(0)
	maximumUncleCount := uint(0)
	posdaoTransition := uint64(0)

	genesis.Config.Aura = &params.AuRaConfig{
		StepDuration: &stepDuration,
		BlockReward:  &blockReward,
		Validators: &params.ValidatorSetJson{
			Multi: map[uint64]*params.ValidatorSetJson{
				0: {List: []common.Address{validatorAddr}},
			},
		},
		BlockRewardContractAddress:    &blockRewardAddr,
		BlockRewardContractTransition: &blockRewardContractTransition,
		MaximumUncleCountTransition:   &maximumUncleCountTransition,
		MaximumUncleCount:             &maximumUncleCount,
		RandomnessContractAddress: map[uint64]common.Address{
			0: randomnessAddr,
		},
		BlockGasLimitContractTransitions: map[uint64]common.Address{
			0: blockGasLimitAddr,
		},
		WithdrawalContractAddress: &withdrawalAddr,
		Registrar:                 &registrarAddr,
		PosdaoTransition:          &posdaoTransition,
		Eip1559FeeCollector:       &feeCollectorAddr,
	}

	return nil
}
