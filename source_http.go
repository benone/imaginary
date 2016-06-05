package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const ImageSourceTypeHttp ImageSourceType = "http"

type HttpImageSource struct {
	Config *SourceConfig
}

func NewHttpImageSource(config *SourceConfig) ImageSource {
	return &HttpImageSource{config}
}

func (s *HttpImageSource) Matches(r *http.Request) bool {
	return r.Method == "GET" && r.URL.Query().Get("url") != ""
}

func (s *HttpImageSource) GetImage(req *http.Request) ([]byte, error) {
	url, err := parseURL(req)
	if err != nil {
		return nil, ErrInvalidImageURL
	}
	if shouldRestrictOrigin(url, s.Config.AllowedOrigings) {
		return nil, fmt.Errorf("Not allowed remote URL origin: %s", url.Host)
	}
	return s.fetchImage(url)
}

func (s *HttpImageSource) fetchImage(url *url.URL) ([]byte, error) {
	req := newHTTPRequest(url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image: (status=%d) (url=%s)", res.StatusCode, req.URL.String())
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %s (url=%s)", req.URL.String(), err)
	}
	return buf, nil
}

func parseURL(request *http.Request) (*url.URL, error) {
  
	queryUrl := request.URL.Query().Get("url")
	sDec, _ := b64.StdEncoding.DecodeString(queryUrl)
  fmt.Println(string(sDec))
  return url.Parse(string(sDec))
}

func newHTTPRequest(url *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("User-Agent", "imaginary/"+Version)
	req.URL = url
	return req
}

func shouldRestrictOrigin(url *url.URL, origins []*url.URL) bool {
	if len(origins) == 0 {
		return false
	}
	for _, origin := range origins {
		if origin.Host == url.Host {
			return false
		}
	}
	return true
}

func init() {
	RegisterSource(ImageSourceTypeHttp, NewHttpImageSource)
}
