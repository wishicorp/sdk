// httputil micro
package httputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
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
	SetHeader(h map[string]string)
	Get(reqUrl string, mediaType MediaType) (ret []byte, err error)
	Put(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error)
	Post(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error)
	Delete(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error)
}

type client struct {
	sync.Mutex
	Timeout time.Duration
	Headers map[string]string
	client  *http.Client
}

func DefaultClient() Client {
	cli := &http.Client{Timeout: DefaultHttpTimeout}
	defer cli.CloseIdleConnections()
	return &client{client: cli, Headers: make(map[string]string)}
}

func NewClient(transport *http.Transport, timeout time.Duration, headers map[string]string) (Client, error) {
	cli := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	defer cli.CloseIdleConnections()

	return &client{
		Headers: headers,
		client:  cli,
	}, nil
}
func NewHttpsClient(certPEMBlock, keyPEMBlock []byte, timeout time.Duration, headers map[string]string) (Client, error) {
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

	return NewClient(transport, timeout, headers)
}

//设置http头
func (c *client) SetHeader(h map[string]string) {
	c.Headers = h
}

func (c *client) Get(reqUrl string, mediaType MediaType) (ret []byte, err error) {
	request, err := http.NewRequest("GET", reqUrl, nil)
	if nil != err {
		return nil, err
	}

	if "" != mediaType {
		request.Header.Add("Accept", string(mediaType))
	}
	return doRequest(c, request)
}

func (c *client) Put(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error) {
	return c.write(reqUrl, PUT, mediaType, data)
}

func (c *client) Post(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error) {
	return c.write(reqUrl, POST, mediaType, data)
}

func (c *client) Delete(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error) {
	return c.write(reqUrl, DELETE, mediaType, data)
}

func (c *client) write(reqUrl string, method RequestMethod, mediaType MediaType, data interface{}) (ret []byte, err error) {

	c.Headers["Content-Kind"] = string(mediaType)
	reader, err := parseReader(mediaType, data)
	if nil != err {
		return nil, err
	}
	request, err := http.NewRequest(string(method), reqUrl, reader)
	if err != nil {
		return nil, err
	}
	return doRequest(c, request)
}

func doRequest(c *client, request *http.Request) (ret []byte, err error) {
	if nil != c.Headers {
		for k, v := range c.Headers {
			request.Header.Set(k, v)
		}
	}

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

func parseReader(mediaType MediaType, data interface{}) (io.Reader, error) {

	switch data.(type) {
	case []byte:
		return bytes.NewReader(data.([]byte)), nil
	case string:
		return strings.NewReader(data.(string)), nil
	default:
		return toMediaTypeReader(mediaType, data)
	}
}

func toMediaTypeReader(mediaType MediaType, data interface{}) (io.Reader, error) {

	switch mediaType {
	case JSON, JsonUtf8:
		reqBytes, err := json.Marshal(data)
		if nil != err {
			return nil, err
		}
		return bytes.NewReader(reqBytes), nil
	case FORM:
		reqBytes, err := json.Marshal(data)
		if nil != err {
			return nil, err
		}
		m := make(map[string]string)
		err = json.Unmarshal(reqBytes, &m)
		if nil != err {
			return nil, err
		}
		u := url.Values{}
		for k, v := range m {
			u.Add(k, v)
		}
		return strings.NewReader(u.Encode()), nil

	case XML, XmlUtf8:
		reqBytes, err := xml.Marshal(data)
		if nil != err {
			return nil, err
		}
		fmt.Println(string(reqBytes))
		return bytes.NewReader(reqBytes), nil
	default:
		return nil, errors.New("media type error")
	}
}

func Get(url string, mediaType MediaType) (ret []byte, err error) {
	c := DefaultClient()
	return c.Get(url, mediaType)
}

func Put(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error) {
	c := DefaultClient()
	return c.Put(reqUrl, mediaType, data)
}
func Post(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error) {
	c := DefaultClient()
	return c.Post(reqUrl, mediaType, data)
}

func Delete(reqUrl string, mediaType MediaType, data interface{}) (ret []byte, err error) {
	c := DefaultClient()
	return c.Delete(reqUrl, mediaType, data)
}

func addTrust(pool *x509.CertPool, path string) error {
	aCrt, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	pool.AppendCertsFromPEM(aCrt)
	return nil
}
