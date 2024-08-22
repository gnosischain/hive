package tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGetValueFromPlaceholders(t *testing.T) {

	responses := []map[string]interface{}{
		{
			"some_value": "Hello, world!",
			"other_key":  "Another value",
		},
		{
			"some_value": []interface{}{
				map[string]interface{}{"some_another_value": "Hello, nested world!"},
			},
			"other_key": "Yet another value",
		},
	}
	// Example query
	query := `$responses[0]["some_value"]`

	// Retrieve the value
	value := getPlaceholderFromContext(responses, query, "responses")

	expected := `Hello, world!`
	if value.(string) != expected {
		t.Errorf("Expected %s, got %s", expected, value.(string))
	}
	query = `$responses[1]["some_value"][0]["some_another_value"]`

	// Retrieve the value
	value = getPlaceholderFromContext(responses, query, "responses")

	expected = `Hello, nested world!`
	if value.(string) != expected {
		t.Errorf("Expected %s, got %s", expected, value.(string))
	}

	query = `$responses[0]["other_key"]`

	// Retrieve the value
	value = getPlaceholderFromContext(responses, query, "responses")

	expected = `Another value`
	if value.(string) != expected {
		t.Errorf("Expected %s, got %s", expected, value.(string))
	}
}

func TestSetJsonValueSimple(t *testing.T) {

	jsonMap := make(map[string]interface{})

	setJsonValue(jsonMap, "jsonrpc", "2.0")
	setJsonValue(jsonMap, "id", "5")
	setJsonValue(jsonMap, "method", "engine_forkchoiceUpdatedV2")
	setJsonValue(jsonMap, "params.[0].headBlockHash", "0x7c6f2d58e5b5cebcbe1ed95c87eadcebf8bc7f520fa7d2c4b04fa6f509661f1a")
	setJsonValue(jsonMap, "params.[1].withdrawals.[0].index", "0x10")
	setJsonValue(jsonMap, "params.[1].withdrawals.[0].address", "0x0000000000000000000000000000000000000000")
	setJsonValue(jsonMap, "params.[1].parentBeaconBlockRoot", "None")
	// convert jsonMap to json string
	jsonStr, _ := json.Marshal(jsonMap)
	fmt.Println(string(jsonStr))

	expected := `{"id":"5","jsonrpc":"2.0","method":"engine_forkchoiceUpdatedV2","params":[{"headBlockHash":"0x7c6f2d58e5b5cebcbe1ed95c87eadcebf8bc7f520fa7d2c4b04fa6f509661f1a"},{"parentBeaconBlockRoot":null,"withdrawals":[{"address":"0x0000000000000000000000000000000000000000","index":"0x10"}]}]}`
	if string(jsonStr) != expected {
		t.Errorf("Expected %s, got %s", expected, string(jsonStr))
	}
}

func TestSetJsonValue(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		value    interface{}
		expected map[string]interface{}
	}{
		{
			name:  "Single value",
			path:  "params.[0].headBlockHash",
			value: "0x7c6f2d58e5b5cebcbe1ed95c87eadcebf8bc7f520fa7d2c4b04fa6f509661f1a",
			expected: map[string]interface{}{
				"params": []interface{}{
					map[string]interface{}{
						"headBlockHash": "0x7c6f2d58e5b5cebcbe1ed95c87eadcebf8bc7f520fa7d2c4b04fa6f509661f1a",
					},
				},
			},
		},
		{
			name:  "Nested value in array",
			path:  "params.[1].withdrawals.[0].index",
			value: "0x10",
			expected: map[string]interface{}{
				"params": []interface{}{
					nil,
					map[string]interface{}{
						"withdrawals": []interface{}{
							map[string]interface{}{
								"index": "0x10",
							},
						},
					},
				},
			},
		},
		{
			name:  "None value",
			path:  "params.[1].parentBeaconBlockRoot",
			value: "None",
			expected: map[string]interface{}{
				"params": []interface{}{
					nil,
					map[string]interface{}{
						"parentBeaconBlockRoot": nil,
					},
				},
			},
		},
		{
			name:  "Multiple nested values",
			path:  "params.[1].withdrawals.[0].address",
			value: "0x0000000000000000000000000000000000000000",
			expected: map[string]interface{}{
				"params": []interface{}{
					nil,
					map[string]interface{}{
						"withdrawals": []interface{}{
							map[string]interface{}{
								"address": "0x0000000000000000000000000000000000000000",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonMap := make(map[string]interface{})
			err := setJsonValue(jsonMap, tt.path, tt.value)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(jsonMap, tt.expected) {
				t.Errorf("Got %v, expected %v", jsonMap, tt.expected)
			}
		})
	}
}
