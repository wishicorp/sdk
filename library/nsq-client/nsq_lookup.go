package nsq_client

import (
	"github.com/nsqio/go-nsq"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type lookup struct {
	config *nsq.Config
	topic  string

	mtx sync.RWMutex

	lookupdHTTPAddrs  []string
	lookupdQueryIndex int
}

func NewLookup(config *nsq.Config, lookupdAddress []string, topic string) *lookup {
	return &lookup{lookupdHTTPAddrs: lookupdAddress, config: config, topic: topic}
}

// initiate a connection to any new producers that are identified.
func (r *lookup) LookupNsqdAddress() []string {
	retries := 0

retry:
	endpoint := r.nextLookup()
	var data lookupResp
	err := apiRequestNegotiateV1("GET", endpoint, nil, &data)
	if err != nil {
		retries++
		if retries < 3 {
			goto retry
		}
		return []string{}
	}

	var nsqdAddrs []string
	for _, producer := range data.Producers {
		broadcastAddress := producer.BroadcastAddress
		port := producer.TCPPort
		joined := net.JoinHostPort(broadcastAddress, strconv.Itoa(port))
		nsqdAddrs = append(nsqdAddrs, joined)
	}
	return nsqdAddrs
}

// return the next lookupd endpoint to query
// keeping track of which one was last used
func (r *lookup) nextLookup() string {
	r.mtx.RLock()
	if r.lookupdQueryIndex >= len(r.lookupdHTTPAddrs) {
		r.lookupdQueryIndex = 0
	}
	addr := r.lookupdHTTPAddrs[r.lookupdQueryIndex]
	num := len(r.lookupdHTTPAddrs)
	r.mtx.RUnlock()
	r.lookupdQueryIndex = (r.lookupdQueryIndex + 1) % num

	urlString := addr
	if !strings.Contains(urlString, "://") {
		urlString = "http://" + addr
	}

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	if u.Path == "/" || u.Path == "" {
		u.Path = "/lookup"
	}

	v, err := url.ParseQuery(u.RawQuery)
	v.Add("topic", r.topic)
	u.RawQuery = v.Encode()
	return u.String()
}

type lookupResp struct {
	Channels  []string    `json:"channels"`
	Producers []*peerInfo `json:"producers"`
	Timestamp int64       `json:"timestamp"`
}

type peerInfo struct {
	RemoteAddress    string `json:"remote_address"`
	Hostname         string `json:"hostname"`
	BroadcastAddress string `json:"broadcast_address"`
	TCPPort          int    `json:"tcp_port"`
	HTTPPort         int    `json:"http_port"`
	Version          string `json:"version"`
}
