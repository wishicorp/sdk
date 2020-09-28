**sdk说明**
* 特别鸣谢 hashicorp.com, 非常喜欢他家的几个产品(有部分代码是直接移植过来的)
* 首先说明本sdk集成了一些三方框架，意在提升开发效率而已
* 第三方库或者sdk依赖
  * gin
  * go-xorm
  * rabbitmq
  * nsq-client
  * rocketmq

**目录结构**
<pre>
sdk/
├── flags 参数处理
├── framework 三方框架
│   ├── gin-http
│   └── worm
├── helper 工具包
│   ├── certutil
│   ├── compressutil
│   ├── errutil
│   ├── errwrap
│   ├── httputil
│   ├── jsonutil
│   ├── kvbuilder
│   ├── parseutil
│   ├── structure
│   ├── strutil
│   ├── threadutil
│   ├── tlsutil
│   └── utils
├── library 三方库
│   ├── consul
│   ├── nsq-client
│   ├── nsq_to_elasticsearch
│   ├── rabbitmq
│   ├── redis
│   └── rocket
├── log 日志组件
├── plugin 基于hcl go-plugin插件的实现
│   ├── constants
│   ├── framework
│   ├── gateway plugin Gateway
│   │   ├── consts
│   │   ├── grpc-gateway
│   │   │   └── proto
│   │   └── http-gateway
│   ├── logical 插件逻辑定义
│   ├── pluginregister 插件注册管理
│   ├── pluginutil 插件工具
│   └── proto 插件协议
├── pool fanout模型的多任务并发池
├── queue 优先级队列
└── version
</pre>

**example**
* 插件server的实现
<pre>
func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "account",
		Level:      hclog.Info,
		JSONFormat: true,
		TimeFormat: time.RFC3339,
	})
	plugin.Serve(&plugin.ServeOpts{
		Factory: account.Factory,
		Logger:  logger,
	})
}
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	//do something ....
}

</pre>