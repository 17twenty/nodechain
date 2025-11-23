package nodechain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SerperSearchTool struct{}

func (t *SerperSearchTool) Name() string { return "web_search" }

type serperRequest struct {
	Query string `json:"q"`
}

type serperResponse struct {
	Organic []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	} `json:"organic"`
	Images []struct {
		Link string `json:"imageUrl"`
	} `json:"images"`
}

func (t *SerperSearchTool) Run(input any) (any, error) {
	query, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("SerperSearchTool: input must be string")
	}

	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SERPER_API_KEY not set")
	}

	reqBody, _ := json.Marshal(serperRequest{Query: query})

	req, _ := http.NewRequest(
		"POST",
		"https://google.serper.dev/search",
		bytes.NewBuffer(reqBody),
	)

	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data serperResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// extract useful URLs
	urls := []string{}
	for _, o := range data.Organic {
		urls = append(urls, o.Link)
	}

	imgs := []string{}
	for _, i := range data.Images {
		imgs = append(imgs, i.Link)
	}

	return map[string]any{
		"urls":   urls,
		"images": imgs,
		"raw":    data,
	}, nil
}
