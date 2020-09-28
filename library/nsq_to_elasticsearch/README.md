将nsq的消息发送到es

* 应用接入方式nsq发送日志的消息格式

```
// Declare a struct for Elasticsearch fields
type ElasticMessage struct {
	Index    string      `json:"index"`    //索引名，未制定将使用service作为索引名
	Service  string      `json:"service"`  //服务名称
	Category string      `json:"category"` //类别
	Pod      string      `json:"pod"`      //k8s的pod名称
	Xid      string      `json:"xid"`      //全局事务id
	Level    string      `json:"level"`    //级别
	Time     *time.Time  `json:"time"`     //日志发送时间
	Body     interface{} `json:"body"`     //日志消息
}

```

* 应用启动
``
nsq_to_elasticsearch --topics jdsh-pay,jdsh-message --channel logging  --lookupd-http-address=127.0.0.1:4161 --elasticsearch-addrs http://127.0.0.1:9200
``
