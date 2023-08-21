package http

import (
	"github.com/kitdine/gbase/log"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConns:        4096,
			MaxIdleConnsPerHost: 2048,
			MaxConnsPerHost:     2048,
			IdleConnTimeout:     30 * time.Second,
		},
		//CheckRedirect: nil,
		Timeout: 30 * time.Second,
	}
	logger = log.Named("http")
)

func Get(url string, headers map[string]string) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Named("get").Error("create request failed", zap.String("url", url), zap.Error(err))
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

}

func Post() {

}

func Put() {

}

func Delete() {

}

func Option() {

}

func Head() {

}

func doRequest(request *http.Request) (*int, []byte, error) {
	resp, err := client.Do(request)
	if err != nil {
		logger.Named("request").Error("request call failed", zap.String("url", request.RequestURI), zap.Error(err))
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &resp.StatusCode, nil, err
	}
	return &resp.StatusCode, body, nil
}

func CustomClientSettings(customClient http.Client) {
	client = &customClient
}
