package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type JsonResponse struct {
	Commit struct {
		Totals struct {
			NotC int         `json:"C"`
			B    int         `json:"b"`
			D    int         `json:"d"`
			F    int         `json:"f"`
			H    int         `json:"h"`
			NotM int         `json:"M"`
			C    string      `json:"c"`
			NotN int         `json:"N"`
			P    int         `json:"p"`
			M    int         `json:"m"`
			Diff interface{} `json:"diff"`
			S    int         `json:"s"`
			N    int         `json:"n"`
		} `json:"totals"`
	} `json:"commit"`
}

func main() {

	var repos []string

	token := os.Getenv("CODECOV_TOKEN")
	host := os.Getenv("CODECOV_HOST")
	path := os.Getenv("CODECOV_PATH")
	branch := os.Getenv("CODECOV_BRANCH")

	file, err := os.Open("repos.txt")
	if err != nil {
		log.Fatal("Open", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repos = append(repos, scanner.Text())
		//		printSlice(repos)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {

		url := fmt.Sprintf("https://%s/%s/%s/branch/%s?access_token=%s", host, path, r, branch, token)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal("NewRequest: ", err)
			return
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("ReadAll:", err)
		}

		defer resp.Body.Close()

		var doc JsonResponse
		jsonErr := json.Unmarshal(body, &doc)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		coverage := "Not Reporting"
		if doc.Commit.Totals.C != "" {
			parsedFloat, err := strconv.ParseFloat(doc.Commit.Totals.C, 32)
			coverage = fmt.Sprintf("%.2f", parsedFloat)
			if err != nil {
				log.Fatal("ParseFloat:", err)
			}

		}

		fmt.Printf("- %s : %s\n", r, coverage)
	}
}

func printSlice(s []string) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
