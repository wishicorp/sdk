// This is an NSQ client that reads the specified topic/channel
// and performs HTTP requests (GET/POST) to the specified endpoints

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/bitly/go-hostpool"
	"github.com/bitly/timer_metrics"
	"github.com/nsqio/go-nsq"
)

const Binary = "1.2.0"

func String(app string) string {
	return fmt.Sprintf("%s v%s (built w/%s)", app, Binary, runtime.Version())
}

type StringArray []string

func (a *StringArray) Get() interface{} { return []string(*a) }

func (a *StringArray) Set(s string) error {
	*a = append(*a, s)
	return nil
}

func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}

const (
	ModeAll = iota
	ModeRoundRobin
	ModeHostPool
)

var (
	showVersion = flag.Bool("version", false, "print version string")

	topics      = flag.String("topics", "", "nsq topic (using comma split)")
	channel     = flag.String("channel", "nsq_to_elasticsearch", "nsq channel")
	maxInFlight = flag.Int("max-in-flight", 200, "max number of messages to allow in flight")

	numPublishers = flag.Int("n", 100, "number of concurrent publishers")
	mode          = flag.String("mode", "hostpool", "the upstream request mode options: round-robin, hostpool (default), epsilon-greedy")
	sample        = flag.Float64("sample", 1.0, "% of messages to publish (float b/w 0 -> 1)")
	statusEvery   = flag.Int("status-every", 250, "the # of requests between logging status (per handler), 0 disables")

	elasticsearchAddrs = StringArray{}
	nsqdTCPAddrs       = StringArray{}
	lookupdHTTPAddrs   = StringArray{}
)

func init() {
	flag.Var(&elasticsearchAddrs, "elasticsearch-addrs", "elasticsearch addrs (may be given multiple times)")
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
}

type PublishHandler struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	counter uint64

	Publisher
	addresses StringArray
	mode      int
	hostPool  hostpool.HostPool

	perAddressStatus map[string]*timer_metrics.TimerMetrics
	timermetrics     *timer_metrics.TimerMetrics
}

func (ph *PublishHandler) HandleMessage(m *nsq.Message) error {
	if *sample < 1.0 && rand.Float64() > *sample {
		return nil
	}
	startTime := time.Now()
	switch ph.mode {
	case ModeAll:
		for _, addr := range ph.addresses {
			st := time.Now()
			err := ph.Publish(m.Body)
			if err != nil {
				return err
			}
			ph.perAddressStatus[addr].Status(st)
		}
	case ModeRoundRobin:
		counter := atomic.AddUint64(&ph.counter, 1)
		idx := counter % uint64(len(ph.addresses))
		addr := ph.addresses[idx]
		err := ph.Publish(m.Body)
		if err != nil {
			return err
		}
		ph.perAddressStatus[addr].Status(startTime)
	case ModeHostPool:
		hostPoolResponse := ph.hostPool.Get()
		addr := hostPoolResponse.Host()
		err := ph.Publish(m.Body)
		hostPoolResponse.Mark(err)
		if err != nil {
			return err
		}
		ph.perAddressStatus[addr].Status(startTime)
	}
	ph.timermetrics.Status(startTime)

	return nil
}

// nsq_to_elasticsearch --topic logging --channel logging --lookupd-http-address 127.0.0.1:4161 --elasticsearch-addrs http://127.0.0.1:9200
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var selectedMode int

	cfg := nsq.NewConfig()

	flag.Var(&nsq.ConfigFlag{Config: cfg}, "consumer-opt", "option to passthrough to nsq.Consumer (may be given multiple times, http://godoc.org/github.com/nsqio/go-nsq#Config)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("nsq_to_elasticsearch v%s\n", Binary)
		return
	}

	if *topics == "" || *channel == "" {
		log.Fatal("--topics and --channel are required")
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatal("--nsqd-tcp-address or --lookupd-http-address required")
	}
	if len(nsqdTCPAddrs) > 0 && len(lookupdHTTPAddrs) > 0 {
		log.Fatal("use --nsqd-tcp-address or --lookupd-http-address not both")
	}
	if len(elasticsearchAddrs) == 0 {
		log.Fatal("--elasticsearch-addrs required")
	}
	switch *mode {
	case "round-robin":
		selectedMode = ModeRoundRobin
	case "hostpool", "epsilon-greedy":
		selectedMode = ModeHostPool
	}

	if *sample > 1.0 || *sample < 0.0 {
		log.Fatal("ERROR: --sample must be between 0.0 and 1.0")
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	cfg.UserAgent = fmt.Sprintf("nsq_to_elasticsearch/%s go-nsq/%s", Binary, nsq.VERSION)
	cfg.MaxInFlight = *maxInFlight

	perAddressStatus := make(map[string]*timer_metrics.TimerMetrics)
	if len(elasticsearchAddrs) == 1 {
		// disable since there is only one address
		perAddressStatus[elasticsearchAddrs[0]] = timer_metrics.NewTimerMetrics(0, "")
	} else {
		for _, a := range elasticsearchAddrs {
			perAddressStatus[a] = timer_metrics.NewTimerMetrics(*statusEvery,
				fmt.Sprintf("[%s]:", a))
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	for _, s := range strings.Split(*topics, ",") {

		go worker(ctx, s, *channel+"-"+s, cfg, selectedMode, perAddressStatus)
	}

	for {
		select {
		case <-termChan:
			cancel()
			return
		}
	}

}

func worker(ctx context.Context, topic string, channel string, config *nsq.Config,
	selectedMode int, perAddressStatus map[string]*timer_metrics.TimerMetrics) {
	log.Println("starting topic:", topic, "channel:", channel)

	publisher, err := NewElasticPublisher(elasticsearchAddrs, topic)
	if nil != err {
		log.Fatal("ERROR: create elastic client:", err)
	}

	publisher.Start()

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Fatal(err)
	}

	hostPool := hostpool.New(elasticsearchAddrs)
	if *mode == "epsilon-greedy" {
		hostPool = hostpool.NewEpsilonGreedy(elasticsearchAddrs, 0, &hostpool.LinearEpsilonValueCalculator{})
	}

	handler := &PublishHandler{
		Publisher:        publisher,
		addresses:        elasticsearchAddrs,
		mode:             selectedMode,
		hostPool:         hostPool,
		perAddressStatus: perAddressStatus,
		timermetrics:     timer_metrics.NewTimerMetrics(*statusEvery, "[aggregate]:"),
	}
	consumer.AddConcurrentHandlers(handler, *numPublishers)

	err = consumer.ConnectToNSQDs(nsqdTCPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("stopping...")
			publisher.Shutdown(false)
			consumer.Stop()
			return
		case <-consumer.StopChan:
			publisher.Cleanup()

		}
	}
}
