package kafka

import (
	"crypto/sha512"
	"github.com/IBM/sarama"
	"github.com/golang/glog"
	"github.com/sebasttiano/gobmp/pkg/pub"
	"github.com/xdg-go/scram"
	"log"
	"os"
	"strings"
)

var (
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

type singlePublisher struct {
	topic    string
	config   *sarama.Config
	producer sarama.AsyncProducer
	stopCh   chan struct{}
}

func (s *singlePublisher) PublishMessage(t int, key []byte, msg []byte) error {
	return s.produceMessage(key, msg)
}

func (s *singlePublisher) Stop() {
	close(s.stopCh)
}

func (s *singlePublisher) produceMessage(key []byte, msg []byte) error {
	var k sarama.ByteEncoder
	var m sarama.ByteEncoder
	k = key
	m = msg
	s.producer.Input() <- &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   k,
		Value: m,
	}

	return nil
}

func NewKafkaSinglePublisher(user, pass, topic, brokers string) (pub.Publisher, error) {
	glog.Infof("creating single publisher")

	if glog.V(6) {
		sarama.Logger = log.New(os.Stdout, "[sarama]      ", log.LstdFlags)
	}
	config := sarama.NewConfig()
	config.Version = sarama.V3_3_1_0
	config.ClientID = "gobmp_client"
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Net.SASL.Enable = true
	config.Net.SASL.User = user
	config.Net.SASL.Password = pass
	config.Net.SASL.Handshake = true
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
	}
	config.Net.TLS.Enable = false

	splitBrokers := strings.Split(brokers, ",")
	for _, b := range splitBrokers {
		if err := validator(b); err != nil {
			glog.Errorf("Failed to validate Kafka server address %s with error: %+v", b, err)
			return nil, err
		}
		br := sarama.NewBroker(b)
		if err := waitForBrokerConnection(br, config, brockerConnectTimeout); err != nil {
			glog.Errorf("failed to open connection to the broker with error: %+v\n", err)
			return nil, err
		}
		glog.V(5).Infof("Connected to broker: %s id: %d\n", br.Addr(), br.ID())
	}

	producer, err := sarama.NewAsyncProducer(splitBrokers, config)
	if err != nil {
		glog.Errorf("New Kafka publisher failed to start new async producer with error: %+v", err)
		return nil, err
	}
	glog.V(5).Infof("Initialized Kafka Async producer")

	stopCh := make(chan struct{})
	go readErrors(producer, stopCh)

	return &singlePublisher{
		stopCh:   stopCh,
		topic:    topic,
		config:   config,
		producer: producer,
	}, nil
}

func readErrors(producer sarama.AsyncProducer, stopCh <-chan struct{}) {
	for {
		select {
		case <-producer.Successes():
		case err := <-producer.Errors():
			glog.Errorf("failed to produce message with error: %+v", *err)
		case <-stopCh:
			producer.Close()
			return
		}
	}
}
