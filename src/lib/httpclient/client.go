package httpclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/abyssparanoia/rapid-go/src/lib/log"

	"github.com/davecgh/go-spew/spew"
)

const defaultTimeout time.Duration = 15 * time.Second

// HTTPOption ... add setting to http module
type HTTPOption struct {
	Headers map[string]string
	Timeout time.Duration
}

// Get ... GET request
func Get(ctx context.Context, u string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Warningf(ctx, "create request error: %s", err.Error())
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	return send(ctx, req, opt)
}

// GetForm ... GET request with form values
func GetForm(ctx context.Context, u string, params map[string]string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Warningf(ctx, "create request error: %s", err.Error())
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	req.URL.RawQuery = query.Encode()
	return send(ctx, req, opt)
}

// GetQueryString ... GET request with query string
func GetQueryString(ctx context.Context, u string, qs string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest("GET", u+"?"+qs, nil)
	if err != nil {
		log.Warningf(ctx, "create request error: %s", err.Error())
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	return send(ctx, req, opt)
}

// PostForm ... POST request with form values
func PostForm(ctx context.Context, u string, params map[string]string, opt *HTTPOption) (int, []byte, error) {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	req, err := http.NewRequest("POST", u, strings.NewReader(values.Encode()))
	if err != nil {
		log.Warningf(ctx, "create request error: %s", err.Error())
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	return send(ctx, req, opt)
}

// PostJSON ... POST request with JSON data
func PostJSON(ctx context.Context, url string, json []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		log.Warningf(ctx, "create request error: %s", err.Error())
		return 0, nil, err
	}

	if opt == nil {
		opt = &HTTPOption{}
	}
	for key, value := range opt.Headers {
		req.Header.Set(key, value)
	}
	opt.Headers["Content-Type"] = "application/json"

	return send(ctx, req, opt)
}

// PostBody ... POST request with body data
func PostBody(ctx context.Context, url string, body []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Warningf(ctx, "create request error: %s", err.Error())
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	return send(ctx, req, opt)
}

func send(ctx context.Context, req *http.Request, opt *HTTPOption) (int, []byte, error) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		log.Debugf(ctx, "send http request: %s", dump)
	} else {
		log.Warningf(ctx, "dumb http request error: %s, error=%s", spew.Sdump(req), err.Error())
	}

	client := http.Client{}
	if opt != nil && opt.Timeout > 0 {
		client.Timeout = opt.Timeout
	} else {
		client.Timeout = defaultTimeout
	}

	res, err := client.Do(req)
	if err != nil {
		log.Warningf(ctx, "http request error: %s", err.Error())
		return 0, nil, err
	}

	dump, err = httputil.DumpResponse(res, true)
	if err == nil {
		log.Debugf(ctx, "http response: %s", dump)
	} else {
		log.Warningf(ctx, "dumb http response error: %s, error=%s", spew.Sdump(req), err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Warningf(ctx, "read http response body error: %s, error=%s", spew.Sdump(res), err.Error())
		return res.StatusCode, nil, nil
	}
	defer res.Body.Close()

	return res.StatusCode, body, nil
}
