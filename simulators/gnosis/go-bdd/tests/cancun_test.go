package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/ethereum/hive/simulators/gnosis/go-bdd/config"
)

type rpcFeature struct {
	responseBody string
	lastResponse *http.Response
	lastRequest  *http.Request
	// list of responses
	responses []map[string]interface{}
	requests  []map[string]interface{}
}

func StringContains(s string, substring string) bool {
	return strings.Contains(s, substring)
}

var placeholdersMap = map[string]interface{}{
	"BASE_URL":    "ENV",
	"ENGINE_URL":  "ENV",
	"AUTH_HEADER": "ENV",
	"TEST":        "TEST",
	// Add other functions from the package here if needed
}

var functionMap = map[string]interface{}{
	"StringContaining": StringContains,
	// Add other functions from the package here if needed
}

func gePlaceholder(input string) string {
	re := regexp.MustCompile(`#([a-zA-Z_][a-zA-Z0-9_]*)#`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 2 {
		fmt.Errorf("input string does not match expected pattern")
		return input
	}
	if placeholdersMap[matches[1]].(string) == "ENV" {
		result := os.Getenv(matches[1])
		println(result)
		return result
	}
	return placeholdersMap[matches[1]].(string)
}

func parseFunction(input string) (interface{}, interface{}) {
	re := regexp.MustCompile(`%([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]+)\)%`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		fmt.Errorf("input string does not match expected pattern")
		return nil, nil
	}

	functionName := matches[1]
	argument := matches[2]

	function, exists := functionMap[functionName]
	if !exists {
		fmt.Errorf("function %s not found in package 'is'", functionName)
		return nil, nil
	}
	return function, argument
}

func convertValue(value interface{}) interface{} {
	if strValue, ok := value.(string); ok {
		switch strings.ToLower(strValue) {
		case "false":
			return false
		case "true":
			return true
		case "none":
			return nil
		default:
			return value
		}
	}
	return value
}

// Parse and retrieve the value from the responses array
func getPlaceholderFromContext(responses []map[string]interface{}, query string, context string) interface{} {
	if !strings.HasPrefix(query, "$"+context+"[") {
		return nil
	}
	query = strings.TrimPrefix(query, "$"+context+"[")

	// Extract the initial index for the responses slice
	closingBracketIndex := strings.Index(query, "]")
	if closingBracketIndex == -1 {
		return nil
	}

	indexStr := query[:closingBracketIndex]
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(responses) {
		return nil
	}

	// Move past the closing bracket and any extra characters
	query = query[closingBracketIndex+1:]

	// Check if there is another bracket immediately after the closing bracket
	if len(query) > 0 && query[0] == '[' {
		query = query[1:] // Skip the opening bracket of the next part
	}

	// Recursive function to traverse the nested maps and slices
	var traverse func(data interface{}, path string) interface{}
	traverse = func(data interface{}, path string) interface{} {
		if len(path) == 0 {
			return data
		}

		closingBracketIndex := strings.Index(path, "]")
		if closingBracketIndex == -1 {
			return nil
		}

		indexOrKey := path[:closingBracketIndex]
		nextPath := strings.TrimSpace(path[closingBracketIndex+1:])

		if len(nextPath) > 0 && nextPath[0] == '[' {
			nextPath = nextPath[1:] // Skip the opening bracket of the next part
		}

		// Debug: Print current path and data type
		fmt.Printf("Current path: %s, Data type: %T\n", path, data)

		if indexOrKey[0] >= '0' && indexOrKey[0] <= '9' {
			// Handle array index
			index, err := strconv.Atoi(indexOrKey)
			if err != nil {
				return nil
			}

			slice, ok := data.([]interface{})
			if !ok {
				return nil
			}

			if index < 0 || index >= len(slice) {
				return nil
			}

			return traverse(slice[index], nextPath)

		} else {
			// Handle map key (unquote it if necessary)
			key := strings.Trim(indexOrKey, `"`)
			m, ok := data.(map[string]interface{})
			if !ok {
				return nil
			}

			value, exists := m[key]
			if !exists {
				return nil
			}

			return traverse(value, nextPath)
		}
	}

	// Start the recursive traversal from the selected map
	return traverse(responses[index], query)
}

func setJsonValue(jsonMap map[string]interface{}, path string, value interface{}) error {
	jsonPath := strings.Replace(path, "[", ".", -1)
	jsonPath = strings.Replace(jsonPath, "]", "", -1)
	jsonPath = strings.Replace(jsonPath, "..", ".", -1)

	keys := strings.Split(jsonPath, ".")

	current := jsonMap
	for i, key := range keys {
		if index, err := strconv.Atoi(key); err == nil {
			// This key is an array index
			parentKey := keys[i-1]
			if _, ok := current[parentKey]; !ok {
				current[parentKey] = []interface{}{}
			}
			array := current[parentKey].([]interface{})

			// Ensure the array has enough elements
			if len(array) <= index {
				array = append(array, make([]interface{}, index-len(array)+1)...)
				current[parentKey] = array
			}

			if i == len(keys)-1 {
				// Convert the value to the correct type
				array[index] = convertValue(value)
			} else {
				// Move to the next level in the path
				if array[index] == nil {
					array[index] = make(map[string]interface{})
				}
				current = array[index].(map[string]interface{})
			}

		} else {
			if i == len(keys)-1 {
				// We're at the last key, set the value
				current[key] = convertValue(value)
			} else {
				if !strings.Contains(jsonPath, ".") {
					// If the key doesn't exist yet, create it as a map
					if _, ok := current[key]; !ok {
						current[key] = make(map[string]interface{})
					}
					current = current[key].(map[string]interface{})
				}
			}
		}
	}

	return nil
}

func (c *rpcFeature) parseAndCreateRequest(arg1 *godog.Table) error {
	//input := ""
	lines := make([]string, len(arg1.Rows))
	for i, row := range arg1.Rows[0:] {
		var cells []string
		for _, cell := range row.Cells {
			cells = append(cells, cell.Value)
		}
		lines[i] = strings.Join(cells, "|")
	}
	// lines := strings.Split(input, "\n")
	var method, url string
	headers := make(map[string]string)
	jsonBody := make(map[string]interface{})

	for _, line := range lines {
		fields := strings.Split(strings.Fields(line)[0], "|")
		placeholder := gePlaceholder(fields[2])
		value := getPlaceholderFromContext(c.requests, placeholder, "requests")

		// Safely check and cast to string
		if strValue, ok := value.(string); ok {
			fmt.Printf("Value as string: %s\n", strValue)
			placeholder = strValue
		}
		value = getPlaceholderFromContext(c.responses, placeholder, "responses")

		// Safely check and cast to string
		if strValue, ok := value.(string); ok {
			fmt.Printf("Value as string: %s\n", strValue)
			placeholder = strValue
		}

		if len(fields) < 3 {
			continue
		}
		switch fields[0] {
		case "Method":
			method = fields[1]
			url = placeholder
		case "Headers":
			headers[fields[1]] = placeholder
		case "Json":
			setJsonValue(jsonBody, fields[1], placeholder)
		}
	}

	// Convert jsonBody map to JSON bytes
	jsonBytes, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fmt.Println("Request headers:", req.Header)
	fmt.Println("Raw body:", string(jsonBytes))
	fmt.Println("Request body:", req.Body)
	// Send the request
	var client *http.Client
	if os.Getenv("HIVE_DEBUG") != "" {
		proxyURL, err := req.URL.Parse(os.Getenv("HTTP_PROXY"))
		if err != nil {
			return err
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	// Update the context
	c.lastResponse = resp
	c.responseBody = string(body)
	c.lastRequest = req

	var result map[string]interface{}
	// Convert JSON string to map[string]interface{}
	json.Unmarshal(jsonBytes, &result)
	c.requests = append(c.requests, result)
	json.Unmarshal(body, &result)
	c.responses = append(c.responses, result)
	//return req, nil
	return nil
}

func getValueFromJSONPath(data map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")

	var current interface{} = data
	for _, key := range keys {
		// Assert that the current value is a map
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid path or non-existent key: %s", path)
		}

		current, ok = currentMap[key]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}

	return current, nil
}

func callFunction(fn interface{}, args []reflect.Value) bool {

	// Call the function using reflection
	results := reflect.ValueOf(fn).Call(args)

	// Check if there is one return value and it's a boolean
	if len(results) == 1 && results[0].Kind() == reflect.Bool {
		result := results[0].Bool() // Convert to bool
		return result
	} else {
		fmt.Println("Unexpected return value")
	}
	return false
}

func (c *rpcFeature) compareJSONPathValues(responseBody string, expectations map[string]string) error {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(responseBody), &data)
	if err != nil {
		return err
	}

	for path, expectedValue := range expectations {
		actualValue, err := getValueFromJSONPath(data, path)
		if err != nil {
			return fmt.Errorf("error extracting value for path %s: %v", path, err)
		}
		fn, arg := parseFunction(expectedValue)
		if fn != nil {
			// Use reflection to invoke the function
			args := []reflect.Value{
				reflect.ValueOf(actualValue),
				reflect.ValueOf(arg),
			}
			if callFunction(fn, args) != true {
				return fmt.Errorf("value mismatch for path %s: function %s, expected %s, got %v", path, fn, reflect.ValueOf(arg), actualValue)
			}
		} else {
			actualValueStr := fmt.Sprintf("%v", actualValue)
			expectedValueStr := fmt.Sprintf("%v", convertValue(expectedValue))
			expectedValueStr = gePlaceholder(expectedValueStr)
			value := getPlaceholderFromContext(c.requests, expectedValueStr, "requests")

			// Safely check and cast to string
			if strValue, ok := value.(string); ok {
				fmt.Printf("Value as string: %s\n", strValue)
				expectedValueStr = strValue
			}
			value = getPlaceholderFromContext(c.responses, expectedValueStr, "responses")

			// Safely check and cast to string
			if strValue, ok := value.(string); ok {
				fmt.Printf("Value as string: %s\n", strValue)
				expectedValueStr = strValue
			}
			if actualValueStr != expectedValueStr {
				return fmt.Errorf("value mismatch for path %s: expected %s, got %s", path, expectedValueStr, actualValueStr)
			}
		}
	}

	return nil
}

func (c *rpcFeature) iShouldGetJsonResponseWithFollowingProperties(arg1 *godog.Table) error {
	lines := make(map[string]string)
	for _, row := range arg1.Rows[0:] {
		var cells []string
		for _, cell := range row.Cells {
			cells = append(cells, cell.Value)
		}
		lines[cells[0]] = cells[1]
	}
	delete(lines, "Path")
	return c.compareJSONPathValues(c.responseBody, lines)
}

func (c *rpcFeature) iShouldReceiveAResponseWithTheStatus(statusCode int) error {
	if c.lastResponse.StatusCode != statusCode {
		return fmt.Errorf("status code expected %d, but found %d", statusCode, c.lastResponse.StatusCode)
	}
	return nil
}

func (c *rpcFeature) theHeaderShouldBe(arg1, arg2 string) error {
	if c.lastResponse.Header.Get(arg1) != arg2 {
		return fmt.Errorf("header '%s' expected to be '%s', but actual %s", arg1, arg2, c.lastResponse.Header.Get(arg1))
	}
	return nil
}

func theResultShouldEqual(arg1 int) error {
	return godog.ErrPending
}

func TestCancunFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeCancunScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t, // Testing instance that will run subtests.
			Tags:     "@cancun",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeCancunScenario(ctx *godog.ScenarioContext) {
	cancun := &rpcFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		config.PLACEHOLDERS["BASE_URL"] = "http://192.168.3.49:8545/"
		cancun.responseBody = ""
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		fmt.Print(cancun.responseBody)
		return ctx, nil
	})
	ctx.Step(`^I send a request with following params$`, cancun.parseAndCreateRequest)
	ctx.Step(`^I should get json response with following properties:$`, cancun.iShouldGetJsonResponseWithFollowingProperties)
	ctx.Step(`^I should receive a response with the status "([^"]*)"$`, cancun.iShouldReceiveAResponseWithTheStatus)
	ctx.Step(`^the header "([^"]*)" should be "([^"]*)"$`, cancun.theHeaderShouldBe)
	ctx.Step(`^the result should equal (\d+)$`, theResultShouldEqual)
}
