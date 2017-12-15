package utils

import (
	"io/ioutil"
	"net/http"
)

func FetchUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil && resp.StatusCode != 200 {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}
