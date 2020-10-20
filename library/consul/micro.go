package consul

import (
	"fmt"
	"github.com/wishicorp/sdk/helper/httputil"
)

type MicroHTTPClient interface {
	Get(reqUrl string, headers map[string][]string) (ret []byte, err error)
	Put(reqUrl string, headers map[string][]string, data interface{}) (ret []byte, err error)
	Post(reqUrl string, headers map[string][]string, data interface{}) (ret []byte, err error)
	Delete(reqUrl string, headers map[string][]string, data interface{}) (ret []byte, err error)
}

type microClient struct {
	url string
	httputil.Client
}

func NewMicroClient(url string, header map[string][]string) *microClient {
	return &microClient{
		url:    url,
		Client: httputil.DefaultClientWithHeader(header),
	}
}

func (c *client) GetMicroHTTPClient(id string, useLan bool, tags string, header map[string][]string) (MicroHTTPClient, error) {
	host, port, err := c.GetServiceAddrPort(id, useLan, tags)
	if nil != err {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &microClient{
		url:    url,
		Client: httputil.DefaultClientWithHeader(header),
	}, nil
}

func (m *microClient) Get(reqUrl string, header map[string][]string) (ret []byte, err error) {
	url := fmt.Sprintf("%s%s", m.url, reqUrl)
	return m.Client.Get(url, header)
}

func (m *microClient) Put(reqUrl string, header map[string][]string, data interface{}) (ret []byte, err error) {
	url := fmt.Sprintf("%s%s", m.url, reqUrl)
	return m.Client.Put(url, header, data)
}

func (m *microClient) Post(reqUrl string, header map[string][]string, data interface{}) (ret []byte, err error) {
	url := fmt.Sprintf("%s%s", m.url, reqUrl)
	return m.Client.Post(url, header, data)
}

func (m *microClient) Delete(reqUrl string, header map[string][]string, data interface{}) (ret []byte, err error) {
	url := fmt.Sprintf("%s%s", m.url, reqUrl)
	return m.Client.Delete(url, header, data)
}
