package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	var allResults []JsonResult
	var cleanResults []JsonResult

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var result JsonResult
		json.Unmarshal(scanner.Bytes(), &result)
		allResults = append(allResults, result)
	}

	for _, result := range allResults {
		if !HasFalsePositivePrefix(allResults, result) {
			cleanResults = append(cleanResults, result)
		}
	}

	resultJsonBytes, _ := json.Marshal(cleanResults)
	fmt.Println(string(resultJsonBytes))
}

func HasFalsePositivePrefix(allResults []JsonResult, curResult JsonResult) bool {
	clusteredResults := map[string][]JsonResult{}
	prefixedResults := map[string][]JsonResult{}

	for _, result := range allResults {
		key := fmt.Sprintf("%d,%d", result.StatusCode, result.ContentWords)
		clusteredResults[key] = append(clusteredResults[key], result)
	}

	for _, resultCluster := range clusteredResults {

		for _, result := range resultCluster {
			if len(string(result.Input["FUZZ"])) < 2 {
				continue
			}
			prefixedResults[string(result.Input["FUZZ"])[0:2]] = append(prefixedResults[string(result.Input["FUZZ"])[0:2]], result)
		}
	}

	if len(prefixedResults[string(curResult.Input["FUZZ"])[0:2]]) > 2 {
		for _, r := range prefixedResults[string(curResult.Input["FUZZ"])[0:2]] {
			if len(curResult.Input["FUZZ"]) > len(r.Input["FUZZ"]) {
				return true
			}
		}
	}

	return false
}

type JsonResult struct {
	Input            map[string]string   `json:"input"`
	Position         int                 `json:"position"`
	StatusCode       int64               `json:"status"`
	ContentLength    int64               `json:"length"`
	ContentWords     int64               `json:"words"`
	ContentLines     int64               `json:"lines"`
	ContentType      string              `json:"content-type"`
	RedirectLocation string              `json:"redirectlocation"`
	ScraperData      map[string][]string `json:"scraper"`
	Duration         time.Duration       `json:"duration"`
	ResultFile       string              `json:"resultfile"`
	Url              string              `json:"url"`
	Host             string              `json:"host"`
}
