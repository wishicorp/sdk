// httputil micro
package httputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type RequestMethod string

const (
	GET     RequestMethod = "GET"
	PUT     RequestMethod = "PUT"
	POST    RequestMethod = "POST"
	DELETE  RequestMethod = "DELETE"
	HEAD    RequestMethod = "HEAD"
	PATCH   RequestMethod = "PATCH"
	OPTIONS RequestMethod = "OPTIONS"
)

type MediaType string

const (
	TEXT     MediaType = "text/plain"
	JSON     MediaType = "application/json"
	JsonUtf8 MediaType = "application/json;charset=utf8"
	XML      MediaType = "application/xml"
	XmlUtf8  MediaType = "application/xml;charset=utf8"
	FORM     MediaType = "application/x-www-form-urlencoded"
)

const DefaultHttpTimeout = 30 * time.Second

type Client interface {
	Get(reqUrl string, header http.Header) (ret []byte, err error)
	Put(reqUrl string, header http.Header, data interface{}) (ret []byte, err error)
	Post(reqUrl string, header http.Header, data interface{}) (ret []byte, err error)
	Delete(reqUrl string, header http.Header, data interface{}) (ret []byte, err error)
}

type client struct {
	sync.Mutex
	Timeout time.Duration
	Header  http.Header
	client  *http.Client
}

func DefaultClient() Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //跳过验证服务端证书
		},
	}
	cli := &http.Client{Timeout: DefaultHttpTimeout}
	cli.Transport = transport
	defer cli.CloseIdleConnections()
	return &client{client: cli, Header: make(http.Header)}
}
func DefaultClientWithHeader(header http.Header) Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //跳过验证服务端证书
		},
	}
	cli := &http.Client{Timeout: DefaultHttpTimeout}
	cli.Transport = transport
	defer cli.CloseIdleConnections()
	return &client{client: cli, Header: header}
}

func NewClient(transport *http.Transport, timeout time.Duration, header http.Header) (Client, error) {
	cli := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	defer cli.CloseIdleConnections()

	return &client{
		Header: header,
		client: cli,
	}, nil
}
func NewHttpsClient(certPEMBlock, keyPEMBlock []byte, timeout time.Duration, header http.Header) (Client, error) {
	pool := x509.NewCertPool()
	cliCrt, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			Certificates:       []tls.Certificate{cliCrt},
			InsecureSkipVerify: true, //跳过验证服务端证书
		},
	}

	return NewClient(transport, timeout, header)
}

func (c *client) mergeHeader(header http.Header) http.Header {
	headers := make(http.Header)
	for k, v := range c.Header {
		headers[k] = v
	}
	if header != nil {
		for k, v := range header {
			headers[k] = v
		}
	}

	return headers
}
func (c *client) Get(reqUrl string, header http.Header) (ret []byte, err error) {
	request, err := http.NewRequest("GET", reqUrl, nil)
	if nil != err {
		return nil, err
	}

	request.Header = c.mergeHeader(header)

	return doRequest(c, request)
}

func (c *client) Put(reqUrl string, header http.Header, data interface{}) (ret []byte, err error) {
	return c.write(reqUrl, PUT, header, data)
}

func (c *client) Post(reqUrl string, header http.Header, data interface{}) (ret []byte, err error) {
	return c.write(reqUrl, POST, header, data)
}

func (c *client) Delete(reqUrl string, header http.Header, data interface{}) (ret []byte, err error) {
	return c.write(reqUrl, DELETE, header, data)
}

func (c *client) write(reqUrl string, method RequestMethod, header http.Header, data interface{}) (ret []byte, err error) {
	if header.Get("Content-Type") == string(FORM) {
		return c.formWrite(reqUrl, method, header, data)
	}

	var reader io.Reader
	switch data.(type) {
	case string:
		reader = strings.NewReader(data.(string))
	case *string:
		reader = strings.NewReader(*data.(*string))
	case []byte:
		reader = bytes.NewReader(data.([]byte))
	case *[]byte:
		reader = bytes.NewReader(*data.(*[]byte))
	default:
		bs, _ := jsonutil.EncodeJSON(data)
		reader = bytes.NewReader(bs)
	}
	request, err := http.NewRequest(string(method), reqUrl, reader)
	if err != nil {
		return nil, err
	}
	c.mergeHeader(header)
	return doRequest(c, request)
}

func (c *client) formWrite(reqUrl string, method RequestMethod, header http.Header, data interface{}) (ret []byte, err error) {
	var form map[string]string
	if err := jsonutil.Swap(data, &form); err != nil {
		return nil, err
	}
	body := url.Values{}
	for key, val := range form {
		body.Set(key, val)
	}
	cli := http.DefaultClient
	resp, err := cli.Post(reqUrl, "application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func doRequest(c *client, request *http.Request) (ret []byte, err error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s %s %s", request.Method, request.URL.String(), http.StatusText(resp.StatusCode))
	}

	return body, nil
}

func Get(url string, header http.Header) (ret []byte, err error) {
	c := DefaultClient()
	return c.Get(url, header)
}

func Put(reqUrl string, header http.Header, data interface{}) (ret []byte, err error) {
	c := DefaultClient()
	return c.Put(reqUrl, header, data)
}
func Post(reqUrl string, header http.Header, data interface{}) (ret []byte, err error) {
	c := DefaultClient()
	return c.Post(reqUrl, header, data)
}

func Delete(reqUrl string, header http.Header, data interface{}) (ret []byte, err error) {
	c := DefaultClient()
	return c.Delete(reqUrl, header, data)
}

func addTrust(pool *x509.CertPool, path string) error {
	aCrt, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	pool.AppendCertsFromPEM(aCrt)
	return nil
}
