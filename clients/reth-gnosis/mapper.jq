# Removes all empty keys and values in input.
def remove_empty:
  . | walk(
    if type == "object" then
      with_entries(
        select(
          .value != null and
          .value != "" and
          .value != [] and
          .key != null and
          .key != ""
        )
      )
    else .
    end
  )
;

# Converts decimal string to number.
def to_int:
  if . == null then . else .|tonumber end
;

# Converts "1" / "0" to boolean.
def to_bool:
  if . == null then . else
    if . == "1" then true else false end
  end
;

# Replace config in input.
. + {
   "seal": {
    "authorityRound": {
      "step": "0x0",
      "signature": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
    }
  },
  "gasLimit": "0x989680",
  "timestamp": 0,
  # "coinbase": "0x0000000000000000000000000000000000000000",
  "baseFeePerGas": "0x3b9aca00",
  # "nonce": "0x0000000000000000",
  "difficulty": "0x01",
    "alloc": ((.alloc|with_entries(.key|="0x"+.)) * {
        "0x59f80ed315477f9f0059D862713A7b082A599217": {
          "balance": "0xc9f2c9cd04674edea40000000"
        },
        "0xB03a86b3126157C039b55E21D378587CcFc04d45": {
          "balance": "0xc9f2c9cd04674edea40000000"
        },
        "0xcC4e00A72d871D6c328BcFE9025AD93d0a26dF51": {
          "balance": "0xc9f2c9cd04674edea40000000"
        },
        "0x0000000000000000000000000000000000000004": {
            "balance": "1",
           
        },
       
        "0x0000000000000000000000000000000000000008": {
            "balance": "0",
        },
        "0x0000000000000000000000000000000000000009": {
            "balance": "0",

        },
       
        "0x0000000000000000000000000000000000000003": {
            "balance": "1",
        },
       
        
        "0x0000000000000000000000000000000000000006": {
            "balance": "0",
        },
       
        "0x0000000000000000000000000000000000000007": {
            "balance": "0",
        },
        "0x0000000000000000000000000000000000000005": {
            "balance": "0",
        },
        "0x0000000000000000000000000000000000000001": {
            "balance": "1",
        },
        "0x0000000000000000000000000000000000000002": {
            "balance": "1",
        }
  }),
  "config": {
    "ChainName": "Gnosis",
    # "ethash": (if env.HIVE_CLIQUE_PERIOD then null else {} end),
    "clique": (if env.HIVE_CLIQUE_PERIOD == null then null else {
      "period": env.HIVE_CLIQUE_PERIOD|to_int,
    } end),
    "chainId": (if env.HIVE_CHAIN_ID == null then 1 else env.HIVE_CHAIN_ID|to_int end),
    "consensus": "aura",
    "homesteadBlock": env.HIVE_FORK_HOMESTEAD|to_int,
    "daoForkBlock": env.HIVE_FORK_DAO_BLOCK|to_int,
    "daoForkSupport": env.HIVE_FORK_DAO_VOTE|to_bool,
    "eip150Block": env.HIVE_FORK_TANGERINE|to_int,
    "eip150Hash": env.HIVE_FORK_TANGERINE_HASH,
    "eip155Block": env.HIVE_FORK_SPURIOUS|to_int,
    "eip158Block": env.HIVE_FORK_SPURIOUS|to_int,
    "byzantiumBlock": env.HIVE_FORK_BYZANTIUM|to_int,
    "constantinopleBlock": env.HIVE_FORK_CONSTANTINOPLE|to_int,
    "petersburgBlock": env.HIVE_FORK_PETERSBURG|to_int,
    "istanbulBlock": env.HIVE_FORK_ISTANBUL|to_int,
    "muirGlacierBlock": env.HIVE_FORK_MUIR_GLACIER|to_int,
    "berlinBlock": env.HIVE_FORK_BERLIN|to_int,
    "londonBlock": env.HIVE_FORK_LONDON|to_int,
    "arrowGlacierBlock": env.HIVE_FORK_ARROW_GLACIER|to_int,
    "grayGlacierBlock": env.HIVE_FORK_GRAY_GLACIER|to_int,
    "mergeNetsplitBlock": env.HIVE_MERGE_BLOCK_ID|to_int,
    "terminalTotalDifficulty": 0,
    "terminalTotalDifficultyPassed": true,
    "shanghaiTime": env.HIVE_SHANGHAI_TIMESTAMP|to_int,
    "cancunTime": env.HIVE_CANCUN_TIMESTAMP|to_int,
    "pragueTime": env.HIVE_PRAGUE_TIMESTAMP|to_int,
    "eip1559FeeCollectorTransition": 0,
    "burntContract": {
      "0": "0x1559000000000000000000000000000000000000"
    },
    "depositContractAddress": "0xbabe2bed00000000000000000000000000000003",
    "minBlobGasPrice": 1000000000,
    "maxBlobGasPerBlock": 262144,
    "targetBlobGasPerBlock": 131072,
    "blobGasPriceUpdateFraction": 1112826,
    "aura": {
      "stepDuration": 5,
      "blockReward": 0,
      "maximumUncleCountTransition": 0,
      "maximumUncleCount": 0,
      "validators": {
        "multi": {
          "0": {
            "list": [
              "0x5cd99ac2f0f8c25a1e670f6bab19d52aad69d875"
            ]
          }
        }
      },
      "blockRewardContractAddress": "0x2000000000000000000000000000000000000001",
      "blockRewardContractTransition": 0,
      "randomnessContractAddress": {
        "0": "0x3000000000000000000000000000000000000001"
      },
      "posdaoTransition": 0,
      "blockGasLimitContractTransitions": {
        "0": "0x4000000000000000000000000000000000000001"
      },
      "registrar": "0x6000000000000000000000000000000000000000",
      "withdrawalContractAddress": "0xbabe2bed00000000000000000000000000000003",
      "twoThirdsMajorityTransition": 0
    }
  }|remove_empty
}