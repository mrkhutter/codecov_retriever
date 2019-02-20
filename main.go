package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	repos := []string{
		"your-repo-1",
		"your-repo-2",
		"your-repo-3",
	}

	token := os.Getenv("CODECOV_TOKEN")
	host := os.Getenv("CODECOV_HOST")
	path := os.Getenv("CODECOV_PATH")
	branch := os.Getenv("CODECOV_BRANCH")

	for _, r := range repos {

		url := fmt.Sprintf("https://%s/%s/%s/branch/%s?access_token=%s", host, path, r, branch, token)
		fmt.Println(url)

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
			coverage = doc.Commit.Totals.C
		}

		fmt.Printf("%s : %s\n", r, coverage)

	}

}
