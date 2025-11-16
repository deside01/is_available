package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func parseBody(r *http.Request) (links []string, err error) {
	var body CheckBody
	json.NewDecoder(r.Body).Decode(&body)

	for _, v := range body.Links {
		if v != "" {
			links = append(links, v)
		}
	}

	if len(links) == 0 {
		return links, errors.New("empty links")
	}

	return links, nil
}

func parseIntBody(r *http.Request) (links []int, err error) {
	var body GetDataBody
	json.NewDecoder(r.Body).Decode(&body)

	for _, v := range body.LinksList {
		if v != 0 {
			links = append(links, v)
		}
	}

	if len(links) == 0 {
		return links, errors.New("empty links")
	}

	return links, nil
}
