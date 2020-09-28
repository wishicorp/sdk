package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	// Import the Elasticsearch library packages
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticClient interface {
	DoIndex(indexName string, doc *ElasticDoc) error
}

type elastic struct {
	*elasticsearch.Config
	*elasticsearch.Client
}

// Declare a struct for Elasticsearch fields
type ElasticDoc struct {
	Id      string                 `json:"id,omitempty"`    //文档ID，未指定使用uuid
	Index   string                 `json:"index,omitempty"` //索引名，未指定将使用topic
	Message map[string]interface{} `json:"message,omitempty"`
}

func NewElasticsearchClient(cfg elasticsearch.Config) (ElasticClient, error) {
	log.Println(GetRoutineID(), "INF start elastic client")
	// Instantiate a new Elasticsearch client object instance
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Elasticsearch connection error: %+v", err)
		return nil, err
	}
	return &elastic{
		Config: &cfg,
		Client: client,
	}, nil
}
func (ne *elastic) Info() {
	// Have the client instance return a response
	res, err := ne.Client.Info()

	// Deserialize the response into a map.
	if err != nil {
		log.Fatalf("client.Info() ERROR: %+v", err)
	} else {
		log.Printf("client response:  %+v", res)
	}
}

func (ne *elastic) DoIndex(indexName string, doc *ElasticDoc) error {

	// Create a context object for the API calls
	body, err := json.Marshal(doc.Message)
	if nil != err {
		return err
	}

	ctx := context.Background()
	var req esapi.Request

	if doc.Id == "" {
		// Instantiate a request object
		req = esapi.IndexRequest{
			Index:      indexName,
			DocumentID: v1UUID(),
			Body:       bytes.NewReader(body),
			Refresh:    "true",
		}
	} else {
		req = esapi.UpdateRequest{
			Index:      indexName,
			DocumentID: doc.Id,
			Body:       bytes.NewReader(body),
		}
	}
	// Return an API response object from request
	res, err := req.Do(ctx, ne.Client)
	if err != nil {
		log.Printf("ERROR: IndexRequest: %s", err)
		return err
	}
	decodeResponse(res)
	return nil
}

func v1UUID() string {
	id, _ := uuid.NewUUID()
	return id.String()
}

func decodeResponse(res *esapi.Response) error {
	defer res.Body.Close()

	if res.IsError() {
		var resMap map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			log.Printf("ERR parsing the response body: %s", err)
			return err
		}
		log.Printf("ERR %s indexing document: %+v", res.Status(),
			resMap["error"].(map[string]interface{})["root_cause"])
		return errors.New(res.Status())
	} else {
		//Deserialize the response into a map.
		//var resMap map[string]interface{}
		//if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		//	log.Printf("Error parsing the response body: %s", err)
		//} else {
		//	//log.Printf("\nIndexRequest() RESPONSE:")
		//	// Print the response status and indexed document version.
		//	fmt.Println("Status:", res.Status())
		//	//fmt.Println("Result:", resMap["result"])
		//	//fmt.Println("Version:", int(resMap["_version"].(float64)))
		//	//fmt.Println("resMap:", resMap)
		//	//fmt.Println("\n")
		//}
	}
	return nil
}
