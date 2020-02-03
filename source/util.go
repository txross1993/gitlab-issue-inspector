package source

import (
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
