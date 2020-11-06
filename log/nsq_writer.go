package log

import (
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"github.com/wishicorp/sdk/library/nsq-client"
	"time"
)

type ElasticDoc struct {
	Id      string                 `json:"id,omitempty"`    //文档ID，未指定使用uuid
	Index   string                 `json:"index,omitempty"` //索引名，未指定将使用topic
	Message map[string]interface{} `json:"message,omitempty"`
}

type nsqWriter struct {
	topic     string
	indexName string
	client    nsq_client.NSQClient
}

func NewNsqWriter(
	lookupdHTTPAddrs []string,
	topic, indexName string) (*nsqWriter, error) {
	config := nsq_client.Config{
		LookupdHTTPAddrs: lookupdHTTPAddrs,
		ReadTimeout: time.Second * 30,
		WriteTimeout: time.Second * 30,
		LookupdPollInterval: time.Second * 5,
	}
	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:       "nsq-logging",
		Level:      hclog.Info,
		JSONFormat: true,
		TimeFormat: time.RFC3339,
	})
	client, err := nsq_client.NewNsqClient(&config, topic, logger)
	if nil != err {
		return nil, err
	}
	if err := client.CreatePublishers(); err != nil {
		return nil, err
	}
	if err := client.StartProducers(true, 1); err != nil {
		return nil, err
	}
	return &nsqWriter{
		topic:     topic,
		client:    client,
		indexName: indexName,
	}, err
}

func (nw *nsqWriter) Write(p []byte) (n int, err error) {
	doc := ElasticDoc{Message: map[string]interface{}{}, Index: nw.indexName}
	if err := json.Unmarshal(p, &doc.Message); err != nil {
		return 0, err
	}
	ts, ok := doc.Message["@timestamp"]
	if ok {
		tm, err := time.ParseInLocation(time.RFC3339, ts.(string), time.Local)
		if nil != err {
			return 0, err
		}
		doc.Message["@timestamp"] = tm
	}
	data, err := jsonutil.EncodeJSON(&doc)
	if nil != err {
		return 0, err
	}
	nw.client.Send(data)
	return len(p), nil
}
