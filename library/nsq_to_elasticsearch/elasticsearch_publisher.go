package main

import (
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"sync"
	"time"
)

type elasticPublisher struct {
	clients []ElasticClient
	inCh    chan []byte
	topic   string
	sync.Mutex
	stopped *bool
}

func NewElasticPublisher(address []string, topic string) (Publisher, error) {
	cfg := elasticsearch.Config{
		Addresses: address,
	}
	clients := make([]ElasticClient, len(address))
	for i := 0; i < len(address); i++ {
		client, err := NewElasticsearchClient(cfg)
		if nil != err {
			return nil, err
		}
		clients[i] = client
	}
	publisher := &elasticPublisher{clients: clients, topic: topic}
	publisher.inCh = make(chan []byte, 4096)
	stopped := false
	publisher.stopped = &stopped
	return publisher, nil

}

func (p *elasticPublisher) Start() {
	for _, client := range p.clients {
		go p.doWork(client, p.inCh)
	}
}

func (p *elasticPublisher) Publish(msg []byte) error {
	if *p.stopped {
		return errors.New("publisher stopped")
	}
	return p.safeSend(msg)
}
func (p *elasticPublisher) Shutdown(immediately bool) {
	p.Lock()
	defer p.Unlock()
	p.stopped = &immediately
}

func (p *elasticPublisher) Cleanup() {
	log.Println("publisher chan closing...")
	close(p.inCh)
}

func (p *elasticPublisher) safeSend(msg []byte) error {
	defer func() {
		if err := recover(); nil != err {
			log.Println("chan closed")
		}
	}()
	select {
	case p.inCh <- msg:
	case <-time.After(time.Second * 5):
		return errors.New("publisher chan is full")
	}
	return nil
}

func (p *elasticPublisher) doWork(client ElasticClient, inCh <-chan []byte) {
	log.Println(GetRoutineID(), "INF start elastic client worker")
	for {
		select {
		case msg, ok := <-inCh:
			if !ok {
				return
			}
			//var load ElasticMessage
			var doc ElasticDoc
			if err := json.Unmarshal(msg, &doc); err != nil {
				log.Println(err.Error(), string(msg))
				continue
			}
			index := p.topic
			if doc.Index != "" {
				index = doc.Index
			}
			if err := client.DoIndex(index, &doc); err != nil {
				log.Println(err, string(msg))
			}
		}
	}
}
