package fetcher

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yalp/jsonpath"
)

const URL = "https://www.caa.ca/wp/wp-admin/admin-ajax.php"

//go:embed locations.json
var locations []byte

func Locations() ([]string, error) {
	mappedLocations := make(map[string]interface{})
	err := json.Unmarshal(locations, &mappedLocations)
	if err != nil {
		return nil, fmt.Errorf("unable to read local locations mapping: %v", err)
	}

	locations := make([]string, 0)
	for location := range mappedLocations {
		locations = append(locations, location)
	}
	return locations, nil
}

func FetchPrice(location string) (string, error) {
	mappedLocations := make(map[string]interface{})
	err := json.Unmarshal(locations, &mappedLocations)
	if err != nil {
		return "", fmt.Errorf("unable to read local locations mapping: %v", err)
	}
	mappedLocation, err := jsonpath.Read(mappedLocations, "$."+location)
	if err != nil {
		return "", fmt.Errorf("invalid location parameter: %v", err)
	}

	client := http.Client{}
	reqBody := strings.NewReader("action=getCitiesForDropdown&caa_dropdown=ONTARIO")
	resp, err := client.Post(URL, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to GET data from URL: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetcher return status code: %v", resp.Status)
	}

	parsedJSON := make(map[string]interface{})
	err = json.Unmarshal(respBody, &parsedJSON)
	if err != nil {
		return "", fmt.Errorf("unable to parse JSON: %v", err)
	}

	priceJsonPath := fmt.Sprintf(`$[today]["%s"]`, mappedLocation.(string))
	curPrice, err := jsonpath.Read(parsedJSON, priceJsonPath)
	if err != nil {
		return "", fmt.Errorf("failed to find price path in JSON: %v %v", err, priceJsonPath)
	}

	trendJsonPath := fmt.Sprintf(`$[arrow]["%s"]`, mappedLocation.(string))
	curTrend, err := jsonpath.Read(parsedJSON, trendJsonPath)
	if err != nil {
		return "", fmt.Errorf("failed to find trend path in JSON: %v %v", err, trendJsonPath)
	}

	trendArrow := ""
	switch curTrend {
	case "up":
		trendArrow = "↑"
	case "down":
		trendArrow = "↓"
	default:
		trendArrow = "="
	}

	return fmt.Sprintf("%v %v %v", curPrice, curTrend, trendArrow), nil
}
