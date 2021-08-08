package common

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"

	"github.com/odwrtw/jackett"
)

type JackettConfig struct {
	Host          string `yaml:"host,omitempty"`
	ApiKey        string `yaml:"api_key,omitempty"`
	BlackholePath string `yaml:"blackhole_path,omitempty"`
}

type JackettResult struct {
	Tracker       string `json:"Tracker"`
	TrackerID     string `json:"TrackerId"`
	CategoryDesc  string `json:"CategoryDesc"`
	Title         string `json:"Title"`
	GUID          string `json:"Guid"`
	Link          string `json:"Link"`
	BlackholeLink string `json:"BlackholeLink"`
	Size          int64  `json:"Size"`
	Seeders       int    `json:"Seeders"`
	Peers         int    `json:"Peers"`
}

type JackettResponse struct {
	Results  []JackettResult   `json:"Results"`
	Indexers []jackett.Indexer `json:"Indexers"`
}

type JackettClient jackett.Client

func NewJackettClient(url, apiKey string) *JackettClient {
	return &JackettClient{url, apiKey}
}

func (c *JackettClient) Search(query string) (JackettResults, error) {
	v := url.Values{}
	v.Add("apikey", c.APIKey)
	v.Add("Query", query)

	url := c.URL + "/api/v2.0/indexers/all/results?" + v.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &JackettResponse{}
	return r.Results, json.NewDecoder(resp.Body).Decode(r)
}

func (r JackettResult) DownloadToBlackhole() error {
	resp, err := http.Get(r.BlackholeLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type JackettResults []JackettResult

func (rs JackettResults) SortBySeeders() []JackettResult {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Seeders > rs[j].Seeders
	})
	return rs
}
