package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	wasmvm "github.com/CosmWasm/wasmvm/types"
)

func DecodeBase64StargateMsgs(
	jsonBz []byte,
) (newSgMsgs []wasmvm.StargateMsg, err error) {
	var data interface{}
	if err := json.Unmarshal(jsonBz, &data); err != nil {
		return []wasmvm.StargateMsg{}, err
	}

	sgMsgs, err := YieldStargateMsgs(jsonBz)
	if err != nil {
		return
	}
	idx := 0
	for sgMsg := range sgMsgs {
		valueStr := string(sgMsg.Value)
		if _, err := json.Marshal(sgMsg.Value); err == nil {
			newSgMsgs = append(newSgMsgs, sgMsg)
		} else if decodedB64, err := base64.StdEncoding.DecodeString(valueStr); err == nil {
			sgMsg.Value = decodedB64
			newSgMsgs = append(newSgMsgs, sgMsg)
		} else {
			return newSgMsgs, fmt.Errorf(
				"parse error: encountered wasmvm.StargateMsg with unexpected format: %s", sgMsg)
		}
		idx++
	}
	return newSgMsgs, nil
}

// TODO: test
// YieldStargateMsgs parses the JSON and sends wasmvm.StargateMsg objects to a channel
func YieldStargateMsgs(jsonBz []byte) (<-chan wasmvm.StargateMsg, error) {
	var data interface{}
	if err := json.Unmarshal(jsonBz, &data); err != nil {
		return nil, err
	}

	ch := make(chan wasmvm.StargateMsg)
	go func() {
		defer close(ch)
		parseStargateMsgChannel(data, ch)
	}()
	return ch, nil
}

func parseStargateMsgChannel(jsonData any, ch chan<- wasmvm.StargateMsg) {
	switch v := jsonData.(type) {
	case map[string]interface{}:
		if typeURL, ok := v["type_url"].(string); ok {
			if value, ok := v["value"].(string); ok {
				ch <- wasmvm.StargateMsg{
					TypeURL: typeURL,
					Value:   []byte(value),
				}
			}
		}
		for _, value := range v {
			parseStargateMsgChannel(value, ch)
		}
	case []interface{}:
		for _, value := range v {
			parseStargateMsgChannel(value, ch)
		}
	}
}

func DecodeBase64Fields(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == "value" {
				if typeURL, ok := v["type_url"].(string); ok && typeURL != "" {
					if encodedValue, ok := value.(string); ok {
						decodedValue, err := base64.StdEncoding.DecodeString(encodedValue)
						if err == nil {
							v[key] = string(decodedValue)
						}
					}
				}
			} else {
				DecodeBase64Fields(value)
			}
		}
	case []interface{}:
		for i := range v {
			DecodeBase64Fields(v[i])
		}
	}
}
