package source

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// helper functions for fetching GitLab data

func performRequest(url string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", os.Getenv("API_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func buildQuery(endpoint string, filters ...string) string {
	baseUrl := os.Getenv("BASE_URL") + endpoint

	for i, f := range filters {
		if f == "" {
			continue
		}
		if i == 0 {
			baseUrl += fmt.Sprintf("?%s", f)
		} else {
			baseUrl += fmt.Sprintf("&%s", f)
		}
	}

	return baseUrl
}
