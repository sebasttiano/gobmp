package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	DstPort    int    `env:"DST_PORT"`
	SrcPort    int    `env:"SRC_PORT"`
	PerfPort   int    `env:"PERF_PORT"`
	KafkaSrv   string `env:"KAFKA_SRV"`
	KafkaTopic string `env:"KAFKA_TOPIC"`
	KafkaUser  string `env:"KAFKA_USER"`
	KafkaPass  string `env:"KAFKA_PASS"`
	NatsSrv    string `env:"NATS_SRV"`
	Intercept  string `env:"INTERCEPT"`
	SplitAF    string `env:"SPLIT_AF"`
	Dump       string `env:"DUMP"`
	File       string `env:"FILE"`
}

func NewConfig() (Config, error) {
	flags := parseCollectorFlags()

	config := Config{}
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if config.DstPort == 0 {
		config.DstPort = flags.DstPort
	}

	if config.SrcPort == 0 {
		config.SrcPort = flags.SrcPort
	}

	if config.PerfPort == 0 {
		config.PerfPort = flags.PerfPort
	}

	if config.KafkaSrv == "" {
		config.KafkaSrv = flags.KafkaSrv
	}

	if config.KafkaTopic == "" {
		config.KafkaTopic = flags.KafkaTopic
	}

	if config.KafkaUser == "" {
		config.KafkaUser = flags.KafkaUser
	}

	if config.KafkaPass == "" {
		config.KafkaPass = flags.KafkaPass
	}

	if config.NatsSrv == "" {
		config.NatsSrv = flags.NatsSrv
	}

	if config.Intercept == "" {
		config.Intercept = flags.Intercept
	}

	if config.SplitAF == "" {
		config.SplitAF = flags.SplitAF
	}

	if config.Dump == "" {
		config.Dump = flags.Dump
	}

	if config.File == "" {
		config.File = flags.File
	}

	return config, nil
}

func parseCollectorFlags() Config {
	srcPort := flag.Int("source-port", 5000, "port exposed to outside")
	dstPort := flag.Int("destination-port", 5050, "port openBMP is listening")
	kafkaSrv := flag.String("kafka-server", "", "URL to access Kafka server")
	kafkaTopic := flag.String("kafka-topic", "", "Kafka topic name")
	kafkaUser := flag.String("kafka-user", "", "Kafka username")
	kafkaPass := flag.String("kafka-pass", "", "Kafka password")
	natsSrv := flag.String("nats-server", "", "URL to access NATS server")
	intercept := flag.String("intercept", "false", "When intercept set \"true\", all incomming BMP messges will be copied to TCP port specified by destination-port, otherwise received BMP messages will be published to Kafka.")
	splitAF := flag.String("split-af", "true", "When set \"true\" (default) ipv4 and ipv6 will be published in separate topics. if set \"false\" the same topic will be used for both address families.")
	perfPort := flag.Int("performance-port", 56767, "port used for performance debugging")
	dump := flag.String("dump", "", "Dump resulting messages to file when \"dump=file\", to standard output when \"dump=console\" or to NATS when \"dump=nats\"")
	file := flag.String("msg-file", "/tmp/messages.json", "Full path anf file name to store messages when \"dump=file\"")

	flag.Parse()

	_ = flag.Set("logtostderr", "true")

	return Config{
		SrcPort:    *srcPort,
		DstPort:    *dstPort,
		KafkaSrv:   *kafkaSrv,
		KafkaTopic: *kafkaTopic,
		KafkaUser:  *kafkaUser,
		KafkaPass:  *kafkaPass,
		NatsSrv:    *natsSrv,
		Intercept:  *intercept,
		SplitAF:    *splitAF,
		PerfPort:   *perfPort,
		Dump:       *dump,
		File:       *file,
	}
}
