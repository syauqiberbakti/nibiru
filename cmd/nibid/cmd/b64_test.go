package cmd_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	wasmvm "github.com/CosmWasm/wasmvm/types"
	"github.com/NibiruChain/nibiru/app"
	"github.com/NibiruChain/nibiru/cmd/nibid/cmd"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

type SuiteB64 struct {
	suite.Suite
}

func TestSuiteB64(t *testing.T) {
	suite.Run(t, new(SuiteB64))
}

func (s *SuiteB64) TestData() []byte {
	return []byte(`{
        "stargate": {
            "type_url": "/cosmos.staking.v1beta1.MsgUndelegate",
            "value": "Cj9uaWJpMTdwOXJ6d25uZnhjanAzMnVuOXVnN3loaHpndGtodmw5amZrc3p0Z3c1dWg2OXdhYzJwZ3N5bjcwbmoSMm5pYml2YWxvcGVyMXdqNWtma25qa3BjNmpkMzByeHRtOHRweGZqZjd4cWx3eDM4YzdwGgwKBXVuaWJpEgMxMTE="
        },
        "another": {
            "type_url": "/cosmos.staking.v1beta1.MsgDelegate",
            "value": "{\"delegator_address\":\"cosmos1eckjje8r8s48kv0pndgtwvehveedlzlnnshl3e\", \"validator_address\":\"cosmos1n6ndsc04xh2hqf506nhvhcggj0qwguf8ks06jj\", \"amount\":{\"denom\":\"unibi\",\"amount\":\"42\"} }"
        }
    }`)
}

func (s *SuiteB64) TestYieldStargateMsgs() {
	jsonBz := s.TestData()
	sgMsgs, err := cmd.YieldStargateMsgs(jsonBz)
	s.NoError(err, sgMsgs)
	msgCount := 0
	for sgMsg := range sgMsgs {
		fmt.Printf("sgMsg: %s\n", sgMsg)
		msgCount++
	}
	s.Equal(msgCount, 2)
}

func (s *SuiteB64) TestB64ParseConcrete() {
	jsonBz := s.TestData()
	sgMsgs, err := cmd.YieldStargateMsgs(jsonBz)
	var sgMsg wasmvm.StargateMsg
	for msg := range sgMsgs {
		sgMsg = msg
		break
	}

	s.Equal(sgMsg.TypeURL, "/cosmos.staking.v1beta1.MsgUndelegate")
	decodedBz, err := base64.StdEncoding.Strict().DecodeString(string(sgMsg.Value))
	s.NoError(err, sgMsg, decodedBz)

	concrete := new(stakingtypes.MsgUndelegate)
	encCfg := app.MakeEncodingConfig()
	// err = encCfg.Marshaler.UnmarshalJSON(decodedBz, concrete)
	err = encCfg.Marshaler.Unmarshal(decodedBz, concrete)
	repr := fmt.Sprintf("decodedBz: %s \nconcrete: %s \n", decodedBz, concrete)
	s.NoError(err, repr)

	outJson, err := encCfg.Marshaler.MarshalJSON(concrete)
	s.NoError(err)
	fmt.Printf("outJson: %s\n", outJson)
	fmt.Printf("repr: %v\n", repr)

	bz, err := json.MarshalIndent(encCfg.InterfaceRegistry.ListImplementations(sdk.MsgInterfaceProtoName), "", "  ")
	s.NoError(err)
	fmt.Printf("bz: %s\n", bz)
}

func (s *SuiteB64) TestCosmosMsgStargate() {
	jsonPayload := `{
        "stargate": {
            "type_url": "/cosmos.staking.v1beta1.MsgUndelegate",
            "value": "Cj9uaWJpMTdwOXJ6d25uZnhjanAzMnVuOXVnN3loaHpndGtodmw5amZrc3p0Z3c1dWg2OXdhYzJwZ3N5bjcwbmoSMm5pYml2YWxvcGVyMXdqNWtma25qa3BjNmpkMzByeHRtOHRweGZqZjd4cWx3eDM4YzdwGgwKBXVuaWJpEgMxMTE="
        },
        "another": {
            "type_url": "/cosmos.staking.v1beta1.MsgDelegate",
            "value": "anothervalue"
        }
    }`

	sgMsgs, err := cmd.DecodeBase64StargateMsgs([]byte(jsonPayload))
	s.NoError(err, "got: ", sgMsgs)

	// encCfg := app.MakeEncodingConfig()

	// cosmosMsg := wasmvm.CosmosMsg{Stargate: asSgMsg}

}
