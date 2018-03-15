package util

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	httpClient http.Client
)

const (
	HTTP_TYPE_GET = iota
	HTTP_TYPE_POST
)

// HTTPServerManager http服务manager
type HTTPServerManager struct {
}

// HTTPServerConfig HTTP服务配置
type HTTPServerConfig struct {
	URL        string   `json:"url"`
	ServerName string   `json:"server_name"`
	Type       int      `json:"type"`
	Request    []string `json:"request"` //KeyName:KeyType
}

// Logger 日志接口
type Logger interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
}

func init() {
	timeout := 10 // GetHTTPTimeout()
	httpClient = http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

// HTTPRequest http common request method
func (p *HTTPServerManager) HTTPRequest(requestParam map[string]string, httpConfig HTTPServerConfig, logger Logger) (respBytes []byte, err error) {
	reqURL := httpConfig.URL
	httpType := "POST"
	if httpConfig.Type == HTTP_TYPE_GET {
		httpType = "GET"
	}

	var req *http.Request

	if httpConfig.Type == HTTP_TYPE_GET {
		req, err = http.NewRequest(httpType, reqURL, nil)
		if err != nil {
			logger.Error("init HTTP GET request url:" + reqURL + " error:" + err.Error())
			return nil, err
		}
		queryValues := req.URL.Query()
		for _, keyName := range httpConfig.Request {
			if value, ok := requestParam[keyName]; ok {
				queryValues.Add(keyName, value)
			} else {
				err := errors.New("Get requset param:" + keyName + " does not exists")
				return nil, err
			}
		}
		req.URL.RawQuery = queryValues.Encode()

		logger.Info("HTTP GET request url:" + req.URL.String())
	} else if httpConfig.Type == HTTP_TYPE_POST {
		form := url.Values{}
		for _, keyName := range httpConfig.Request {
			if value, ok := requestParam[keyName]; ok {
				form.Add(keyName, value)
			} else {
				err := errors.New("POST requset param:" + keyName + " does not exists")
				return nil, err
			}
		}

		req, err = http.NewRequest(httpType, reqURL, strings.NewReader(form.Encode()))
		if err != nil {
			logger.Error("init HTTP POST request url:" + reqURL + " error:" + err.Error())
			return nil, err
		}

		logger.Info("HTTP POST request url:" + reqURL + " param:" + form.Encode())
		req.Header.Add("Content-Type", "application/json")
	} else {
		err := errors.New("not support HTTP method except GET or POST")
		return nil, err
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		logger.Error("HTTP " + httpType + " Request url:" + reqURL + " error:" + err.Error())
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			respBytes = nil
		}
	}()

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Fetch HTTP " + httpType + " Response error:" + err.Error())
		return nil, err
	}
	return respBytes, err
}
