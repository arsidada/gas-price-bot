package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/yalp/jsonpath"
)

const URL = "https://www.caa.ca/wp/wp-admin/admin-ajax.php"

func main() {
	client := http.Client{}
	reqBody := strings.NewReader("action=getCitiesForDropdown&caa_dropdown=ONTARIO")
	resp, err := client.Post(URL, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		log.Fatalf("Failed to GET data from URL: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	parsedJSON := make(map[string]interface{})
	err = json.Unmarshal(respBody, &parsedJSON)
	if err != nil {
		log.Fatalf("Unable to parse JSON: %v", err)
	}

	curPrice, err := jsonpath.Read(parsedJSON, "$.today.HALTON")
	if err != nil {
		log.Fatalf("Failed to find price path in JSON: %v", err)
	}

	curTrend, err := jsonpath.Read(parsedJSON, "$.arrow.HALTON")
	if err != nil {
		log.Fatalf("Failed to find trend path in JSON: %v", err)
	}

	fmt.Printf("Halton gas price: %s ", curPrice)
	switch curTrend {
	case "up":
		fmt.Println("↑")
	case "down":
		fmt.Println("↓")
	default:
		fmt.Println("=")
	}
}
