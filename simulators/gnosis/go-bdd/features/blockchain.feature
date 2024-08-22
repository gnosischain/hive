Feature: Blockchain tests

    @curl @cancun
    Scenario: eth_getBlockByNumber test
        When I send a request with following params
            | Method  | POST         | #BASE_URL#           |
            | Headers | Content-Type | application/json     |
            | Json    | jsonrpc      | 2.0                  |
            | Json    | id           | 1                    |
            | Json    | method       | eth_getBlockByNumber |
            | Json    | params[0]    | latest               |
            | Json    | params[1]    | false                |

        Then I should get json response with following properties:
            | Path                    | Value                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
            | jsonrpc                 | $requests[0]["jsonrpc"]                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
            | result.author           | 0x0000000000000000000000000000000000000000                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
            | result.difficulty       | 0x64                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
            | result.extraData        | 0x                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | result.gasLimit         | 0x989680                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
            | result.gasUsed          | 0x0                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
            | result.hash             | 0xa01dbd45233583bc597c05a970684755b5016664842f8f3d394fff6bc4852404                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | result.logsBloom        | 0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 |
            | result.miner            | 0x0000000000000000000000000000000000000000                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
            | result.number           | 0x0                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
            | result.parentHash       | 0x0000000000000000000000000000000000000000000000000000000000000000                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | result.receiptsRoot     | 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | result.sha3Uncles       | 0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | result.signature        | 0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000                                                                                                                                                                                                                                                                                                                                                                                               |
            | result.size             | 0x219                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
            | result.stateRoot        | 0xd24141462904e0137d827b846bb691d9806d0c0396ffa05ca209003c162afdaa                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | result.totalDifficulty  | %StringContaining(0x64)%                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
            | result.timestamp        | 0x1234                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
            | result.transactionsRoot | 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
            | id                      | 1                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
        And I should receive a response with the status "200"
        And the header "Content-Type" should be "application/json"



    @curl @cancun
    Scenario: engine_forkchoiceUpdatedV2 SYNCING test
        Given I send a request with following params
            | Method  | POST         | #BASE_URL#           |
            | Headers | Content-Type | application/json     |
            | Json    | jsonrpc      | 2.0                  |
            | Json    | id           | 1                    |
            | Json    | method       | eth_getBlockByNumber |
            | Json    | params[0]    | latest               |
            | Json    | params[1]    | true                 |
       When I send a request with following params
            | Method  | POST                                      | #ENGINE_URL#                                                       |
            | Headers | Authorization                             | #AUTH_HEADER#                                                      |
            | Headers | Content-Type                              | application/json                                                   |
            | Json    | jsonrpc                                   | $requests[0]["jsonrpc"]                                            |
            | Json    | id                                        | 2                                                                  |
            | Json    | method                                    | engine_forkchoiceUpdatedV2                                         |
            # Hash does not meet last block hash
            | Json    | params.[0].headBlockHash                  | 0x7c6f2d58e5b5cebcbe1ed95c87eadcebf8bc7f520fa7d2c4b04fa6f509661f1a |
            | Json    | params.[0].safeBlockHash                  | 0x0000000000000000000000000000000000000000000000000000000000000000 |
            | Json    | params.[0].finalizedBlockHash             | 0x0000000000000000000000000000000000000000000000000000000000000000 |
            | Json    | params.[1].timestamp                      | 0x66bf5f8d                                                         |
            | Json    | params.[1].prevRandao                     | 0x1ab713097b6f9d6619115f59da52168b37b179b031ceb909db6a7632274183ea |
            | Json    | params.[1].suggestedFeeRecipient          | 0x0000000000000000000000000000000000000000                         |
            | Json    | params.[1].parentBeaconBlockRoot          | None                                                               |
            | Json    | params.[1].withdrawals.[0].index          | 0x10                                                               |
            | Json    | params.[1].withdrawals.[0].validatorIndex | 0x0                                                                |
            | Json    | params.[1].withdrawals.[0].address        | 0x0000000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[0].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[1].index          | 0x11                                                               |
            | Json    | params.[1].withdrawals.[1].validatorIndex | 0x1                                                                |
            | Json    | params.[1].withdrawals.[1].address        | 0x0100000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[1].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[2].index          | 0x12                                                               |
            | Json    | params.[1].withdrawals.[2].validatorIndex | 0x2                                                                |
            | Json    | params.[1].withdrawals.[2].address        | 0x0200000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[2].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[3].index          | 0x13                                                               |
            | Json    | params.[1].withdrawals.[3].validatorIndex | 0x3                                                                |
            | Json    | params.[1].withdrawals.[3].address        | 0x0300000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[3].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[4].index          | 0x14                                                               |
            | Json    | params.[1].withdrawals.[4].validatorIndex | 0x4                                                                |
            | Json    | params.[1].withdrawals.[4].address        | 0x0400000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[4].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[5].index          | 0x15                                                               |
            | Json    | params.[1].withdrawals.[5].validatorIndex | 0x5                                                                |
            | Json    | params.[1].withdrawals.[5].address        | 0x0500000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[5].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[6].index          | 0x16                                                               |
            | Json    | params.[1].withdrawals.[6].validatorIndex | 0x6                                                                |
            | Json    | params.[1].withdrawals.[6].address        | 0x0600000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[6].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[7].index          | 0x17                                                               |
            | Json    | params.[1].withdrawals.[7].validatorIndex | 0x7                                                                |
            | Json    | params.[1].withdrawals.[7].address        | 0x0700000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[7].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[8].index          | 0x18                                                               |
            | Json    | params.[1].withdrawals.[8].validatorIndex | 0x8                                                                |
            | Json    | params.[1].withdrawals.[8].address        | 0x0800000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[8].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[9].index          | 0x19                                                               |
            | Json    | params.[1].withdrawals.[9].validatorIndex | 0x9                                                                |
            | Json    | params.[1].withdrawals.[9].address        | 0x0900000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[9].amount         | 0x64                                                               |
        Then I should get json response with following properties:
            | Path                                 | Value                    |
            | jsonrpc                              | $responses[1]["jsonrpc"] |
            | result.payloadStatus.status          | SYNCING                  |
            | result.payloadStatus.latestValidHash | None                     |
            | result.payloadStatus.validationError | None                     |
            | result.payloadId                     | None                     |
            | id                                   | 2                        |

    @curl @cancun
    Scenario: Withdrawals test
        When I send a request with following params
            | Method  | POST         | #BASE_URL#           |
            | Headers | Content-Type | application/json     |
            | Json    | jsonrpc      | 2.0                  |
            | Json    | id           | 1                    |
            | Json    | method       | eth_getBlockByNumber |
            | Json    | params[0]    | latest               |
            | Json    | params[1]    | true                 |

        # When I send a request with following params
        #     | Method  | POST         | #BASE_URL#                                                         |
        #     | Headers | Content-Type | application/json                                                   |
        #     | Json    | jsonrpc      | 2.0                                                                |
        #     | Json    | id           | 1                                                                  |
        #     | Json    | method       | eth_getStorageAt                                                   |
        #     | Json    | params[0]    | 0x0000000000000000000000000000000000001003                         |
        #     | Json    | params[1]    | 0x0000000000000000000000000000000000000000000000000000000000000000 |
        #     | Json    | params[1]    | latest                                                             |
        # When I send a request with following params
        #     | Method  | POST         | #BASE_URL#                                                                                                                                                                                                                           |
        #     | Headers | Content-Type | application/json                                                                                                                                                                                                                     |
        #     | Json    | jsonrpc      | 2.0                                                                                                                                                                                                                                  |
        #     | Json    | id           | 2                                                                                                                                                                                                                                    |
        #     | Json    | method       | eth_sendRawTransaction                                                                                                                                                                                                               |
        #     | Json    | params[0]    | 0x02f86e8227da0b843b9aca008506fc23ac00830124f89402020202020202020202020202020202020202020180c001a03f6f64f57a28c950d58cdb60ac7b0aa60a6e743e58f5f8b47b0de9095678ca54a06755185aaca640d9ecdcc3140bfd3e87cb732771fb246bbf785ca5b75910cefb |
        When I send a request with following params
            | Method  | POST                                      | #ENGINE_URL#                                                       |
            | Headers | Authorization                             | #AUTH_HEADER#                                                      |
            | Headers | Content-Type                              | application/json                                                   |
            | Json    | jsonrpc                                   | $requests[0]["jsonrpc"]                                            |
            | Json    | id                                        | 2                                                                  |
            | Json    | method                                    | engine_forkchoiceUpdatedV2                                         |
            | Json    | params.[0].headBlockHash                  | $requests[0]["result"]["hash"]                                     |
            | Json    | params.[0].safeBlockHash                  | 0x0000000000000000000000000000000000000000000000000000000000000000 |
            | Json    | params.[0].finalizedBlockHash             | 0x0000000000000000000000000000000000000000000000000000000000000000 |
            | Json    | params.[1].timestamp                      | 1724134336                                                         |
            | Json    | params.[1].prevRandao                     | 0x1ab713097b6f9d6619115f59da52168b37b179b031ceb909db6a7632274183ea |
            | Json    | params.[1].suggestedFeeRecipient          | 0x0000000000000000000000000000000000000000                         |
            | Json    | params.[1].parentBeaconBlockRoot          | None                                                               |
            | Json    | params.[1].withdrawals.[0].index          | 0x10                                                               |
            | Json    | params.[1].withdrawals.[0].validatorIndex | 0x0                                                                |
            | Json    | params.[1].withdrawals.[0].address        | 0x0000000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[0].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[1].index          | 0x11                                                               |
            | Json    | params.[1].withdrawals.[1].validatorIndex | 0x1                                                                |
            | Json    | params.[1].withdrawals.[1].address        | 0x0100000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[1].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[2].index          | 0x12                                                               |
            | Json    | params.[1].withdrawals.[2].validatorIndex | 0x2                                                                |
            | Json    | params.[1].withdrawals.[2].address        | 0x0200000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[2].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[3].index          | 0x13                                                               |
            | Json    | params.[1].withdrawals.[3].validatorIndex | 0x3                                                                |
            | Json    | params.[1].withdrawals.[3].address        | 0x0300000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[3].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[4].index          | 0x14                                                               |
            | Json    | params.[1].withdrawals.[4].validatorIndex | 0x4                                                                |
            | Json    | params.[1].withdrawals.[4].address        | 0x0400000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[4].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[5].index          | 0x15                                                               |
            | Json    | params.[1].withdrawals.[5].validatorIndex | 0x5                                                                |
            | Json    | params.[1].withdrawals.[5].address        | 0x0500000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[5].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[6].index          | 0x16                                                               |
            | Json    | params.[1].withdrawals.[6].validatorIndex | 0x6                                                                |
            | Json    | params.[1].withdrawals.[6].address        | 0x0600000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[6].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[7].index          | 0x17                                                               |
            | Json    | params.[1].withdrawals.[7].validatorIndex | 0x7                                                                |
            | Json    | params.[1].withdrawals.[7].address        | 0x0700000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[7].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[8].index          | 0x18                                                               |
            | Json    | params.[1].withdrawals.[8].validatorIndex | 0x8                                                                |
            | Json    | params.[1].withdrawals.[8].address        | 0x0800000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[8].amount         | 0x64                                                               |
            | Json    | params.[1].withdrawals.[9].index          | 0x19                                                               |
            | Json    | params.[1].withdrawals.[9].validatorIndex | 0x9                                                                |
            | Json    | params.[1].withdrawals.[9].address        | 0x0900000000000000000000000000000000000000                         |
            | Json    | params.[1].withdrawals.[9].amount         | 0x64                                                               |
        Then I should get json response with following properties:
            | Path                                 | Value                    |
            | jsonrpc                              | $responses[1]["jsonrpc"] |
            | result.payloadStatus.status          | VALID                  |
            #| result.payloadStatus.latestValidHash | None                     |
            #| result.payloadStatus.validationError | None                     |
            #| result.payloadId                     | None                     |
            | id                                   | 2                        |