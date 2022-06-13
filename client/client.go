package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/aexvir/lnk/api"
)

// Lnk is a client for the lnk service.
type Lnk struct {
	baseurl string
	client  *http.Client
}

// NewLnkClient instantiates a new client for the lnk service.
// Not all features are implemented currently.
// Customize it via ClientOpts.
func NewLnkClient(opts ...ClientOpt) (*Lnk, error) {
	client := Lnk{
		baseurl: "http://localhost:8000",
		client: &http.Client{
			Transport: DefaultTransport(),
		},
	}

	for _, opt := range opts {
		err := opt(&client)
		if err != nil {
			return nil, err
		}
	}

	return &client, nil
}

// CreateLink for a target url with an optional custom slug.
func (lc *Lnk) CreateLink(target string, slug *string) (*api.LinkResp, error) {
	req := api.CreateLinkReq{
		Target: target,
		Slug:   slug,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error serializing payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/links", lc.baseurl)
	resp, err := lc.client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed; status: %d", resp.StatusCode)
	}

	var link api.LinkResp
	err = json.NewDecoder(resp.Body).Decode(&link)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &link, nil
}

// GetLink for a specific slug.
func (lc *Lnk) GetLink(slug string) (*api.Link, error) {
	url := fmt.Sprintf("%s/api/links/%s", lc.baseurl, slug)

	resp, err := lc.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var link api.Link
		err = json.NewDecoder(resp.Body).Decode(&link)
		if err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		return &link, nil

	case http.StatusNotFound:
		return nil, fmt.Errorf("slug %s not found", slug)

	default:
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
}

type ClientOpt func(*Lnk) error

func WithBaseUrl(url string) ClientOpt {
	return func(lc *Lnk) error {
		lc.baseurl = url
		return nil
	}
}

func DefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
	}
}
