package api

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const domain = "kequ5060.cn:8200"
const scheme = "http" // https / http
const MaxIdleConns = 3000

var OpenClient = resty.New()

type uri string

type Client struct {
	Client *resty.Client
}

const (
	pricesGetUri uri = "/prices/{service_number}/{keys}"
	priceGetUri  uri = "/prices/{service_number}/{item}"
)

func (c *Client) getUrl(endpoint uri) string {
	return fmt.Sprintf("%s://%s%s", scheme, domain, endpoint)
}

func (c *Client) request(ctx context.Context) *resty.Request {
	return c.Client.R().SetContext(ctx)
}

func GetPricesUrl(ctx context.Context, service_number, keys string) (*resty.Response, error) {
	c := Client{
		Client: resty.New(),
	}
	r, err := c.request(ctx).
		SetPathParams(map[string]string{
			"service_number": service_number,
			"keys":           keys,
		}).
		Patch(c.getUrl(pricesGetUri))
	if err != nil {
		return nil, err
	}
	return r, nil
}

/*func createTransport(localAddr net.Addr, idleConns int) *http.Transport {
	dialer := &net.Dialer{
		Timeout: 60 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	if localAddr != nil {
		dialer.LocalAddr = localAddr
	}
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: dialer.DialContext,
		ForceAttemptHTTP2: true,
		MaxIdleConns: idleConns,
		IdleConnTimeout: 90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost: idleConns,
		MaxConnsPerHost: idleConns,
	}
}*/
