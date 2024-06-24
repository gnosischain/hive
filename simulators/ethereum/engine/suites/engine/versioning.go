package suite_engine

import (
	"github.com/ethereum/hive/simulators/ethereum/engine/clmock"
	"github.com/ethereum/hive/simulators/ethereum/engine/config"
	"github.com/ethereum/hive/simulators/ethereum/engine/helper"
	"github.com/ethereum/hive/simulators/ethereum/engine/test"
)

// Test versioning of the Engine API methods

type EngineNewPayloadVersionTest struct {
	test.BaseSpec
}

func (s EngineNewPayloadVersionTest) WithMainFork(fork config.Fork) test.Spec {
	specCopy := s
	specCopy.MainFork = fork
	return specCopy
}

func (s EngineNewPayloadVersionTest) WithTimestamp(genesisTime uint64) test.Spec {
	specCopy := s
	// Set genesis time if not defined
	if s.GenesisTimestamp == nil {
		specCopy.GenesisTimestamp = &genesisTime
	}
	// Set fork time, will be ignored if fork height is set
	specCopy.ForkTime = *specCopy.GenesisTimestamp
	// Set previous fork time if fork height is set
	mainFork := s.GetMainFork()
	if s.ForkHeight > 0 && mainFork != config.Paris && mainFork != config.Shanghai {
		// No previous fork time for Paris and Shanghai
		specCopy.PreviousForkTime = genesisTime
	}
	return specCopy
}

// Test modifying the ForkchoiceUpdated version on Payload Request to the previous/upcoming version
// when the timestamp payload attribute does not match the upgraded/downgraded version.
type ForkchoiceUpdatedOnPayloadRequestTest struct {
	test.BaseSpec
	helper.ForkchoiceUpdatedCustomizer
}

func (s ForkchoiceUpdatedOnPayloadRequestTest) WithMainFork(fork config.Fork) test.Spec {
	specCopy := s
	specCopy.MainFork = fork
	return specCopy
}

func (s ForkchoiceUpdatedOnPayloadRequestTest) WithTimestamp(genesisTime uint64) test.Spec {
	specCopy := s
	// Set genesis time if not defined
	if s.GenesisTimestamp == nil {
		specCopy.GenesisTimestamp = &genesisTime
	}
	// Set fork time, will be ignored if fork height is set
	specCopy.ForkTime = *specCopy.GenesisTimestamp
	// Set previous fork time if fork height is set
	mainFork := s.GetMainFork()
	if s.ForkHeight > 0 && mainFork != config.Paris && mainFork != config.Shanghai {
		// No previous fork time for Paris and Shanghai
		specCopy.PreviousForkTime = genesisTime
	}
	return specCopy
}

func (tc ForkchoiceUpdatedOnPayloadRequestTest) GetName() string {
	return "ForkchoiceUpdated Version on Payload Request: " + tc.BaseSpec.GetName()
}

func (tc ForkchoiceUpdatedOnPayloadRequestTest) Execute(t *test.Env) {
	// Wait until TTD is reached by this client
	t.CLMock.WaitForTTD()

	t.CLMock.ProduceSingleBlock(clmock.BlockProcessCallbacks{
		OnPayloadAttributesGenerated: func() {
			var (
				payloadAttributes                    = &t.CLMock.LatestPayloadAttributes
				expectedStatus    test.PayloadStatus = test.Valid
				expectedError     *int
				err               error
			)
			tc.SetEngineAPIVersionResolver(t.ForkConfig)
			testEngine := t.TestEngine.WithEngineAPIVersionResolver(tc.ForkchoiceUpdatedCustomizer)
			payloadAttributes, err = tc.GetPayloadAttributes(payloadAttributes)
			if err != nil {
				t.Fatalf("FAIL: Error getting custom payload attributes: %v", err)
			}
			expectedError, err = tc.GetExpectedError()
			if err != nil {
				t.Fatalf("FAIL: Error getting custom expected error: %v", err)
			}
			if tc.GetExpectInvalidStatus() {
				expectedStatus = test.Invalid
			}

			r := testEngine.TestEngineForkchoiceUpdated(&t.CLMock.LatestForkchoice, payloadAttributes, t.CLMock.LatestHeader.Time)
			r.ExpectationDescription = tc.Expectation
			if expectedError != nil {
				r.ExpectErrorCode(*expectedError)
			} else {
				r.ExpectNoError()
				r.ExpectPayloadStatus(expectedStatus)
			}
		},
	})
}
