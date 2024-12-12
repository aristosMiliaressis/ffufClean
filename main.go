package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	if resultJsonBytes != nil {
		fmt.Println(string(resultJsonBytes))
	}
}

func HasFalsePositivePrefix(allResults []JsonResult, curResult JsonResult) bool {
	clusteredResults := map[string][]JsonResult{}
	prefixedResults := map[string][]JsonResult{}

	for _, result := range allResults {
		key := fmt.Sprintf("%d,%d", result.StatusCode, result.ContentWords)
		clusteredResults[key] = append(clusteredResults[key], result)
	}

	key := fmt.Sprintf("%d,%d", curResult.StatusCode, curResult.ContentWords)

	for _, result := range clusteredResults[key] {
		path := strings.SplitN(result.Url, result.Host, 2)[1]
		if len(path) < 2 {
			continue
		}
		prefixedResults[path[0:2]] = append(prefixedResults[path[0:2]], result)
	}

	curPath := strings.SplitN(curResult.Url, curResult.Host, 2)[1]
	if len(prefixedResults[string(curPath)[0:2]]) > 1 {
		for _, r := range prefixedResults[curPath[0:2]] {
			rPath := strings.SplitN(r.Url, r.Host, 2)[1]
			if len(curPath) > len(rPath) {
				fmt.Fprintf(os.Stderr, "%s is false positve", curResult.Url)
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
