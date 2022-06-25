package fw

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
)

type AppHttpClient struct {
	httpClient *http.Client
	maxRetry   int
}

var (
	httpClientOnce sync.Once
	appHttpClient  *AppHttpClient
)

func GetHttpClient() *AppHttpClient {
	httpClientOnce.Do(func() {
		timeOut := 3
		maxRetry := 1
		tr := &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     90 * time.Second,
		}
		appHttpClient = &AppHttpClient{
			httpClient: &http.Client{
				Timeout:   time.Second * time.Duration(timeOut),
				Transport: tr,
			},
			maxRetry: maxRetry,
		}
	})
	return appHttpClient
}

func (httpClient *AppHttpClient) RequestCtx(ctx context.Context, method string, uri string, payload interface{}, response interface{}) (int, error) {
	try := 0
	resp, err := httpClient.request(ctx, method, uri, payload, nil, try, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("rpc response: %d", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return 0, err
		}
	} else if strings.Contains(ct, "application/octet-stream") {
		var respByte []byte
		respByte, err = ioutil.ReadAll(resp.Body)
		pv, ok := response.(*[]byte)
		if !ok {
			return 0, errors.New("unsupported receiver type for application/octet-stream response")
		}
		*pv = respByte
	} else {
		// handle with text
		var respByte []byte
		respByte, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		switch response.(type) {
		case *string:
			pv := response.(*string)
			*pv = string(respByte)
		case *[]byte:
			pv := response.(*[]byte)
			*pv = respByte
		default:
			return 0, errors.New("unsupported receiver type")
		}
	}
	return resp.StatusCode, nil
}

func (httpClient *AppHttpClient) request(ctx context.Context, method, uri string, payload interface{}, header map[string]string, try int, oldErr error) (*http.Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, uri)
	defer span.End()
	if try >= httpClient.maxRetry {
		return nil, oldErr
	}
	if try > 0 {
		time.Sleep(500 * time.Millisecond)
	}
	var body io.Reader
	if payload != nil {
		jsonValue, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonValue)
	}
	var resp *http.Response
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	// set header
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	// content-type
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	err = DoWithCancel(ctx, func() error {
		resp, err = httpClient.httpClient.Do(req.WithContext(ctx))
		return err
	})
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return httpClient.request(ctx, method, uri, payload, header, try+1, err)
	}
	return resp, nil
}

var ErrorTimeout = errors.New("访问超时")

func DoWithCancel(ctx context.Context, f func() error) error {
	done := make(chan error)
	go func() {
		done <- f()
	}()
	select {
	case <-ctx.Done():
		return ErrorTimeout
	case err := <-done:
		return err
	}
}
