package main

import (
	"fmt"
	"github.com/sebasttiano/gobmp/pkg/config"
	"os"
	"runtime"
	"strconv"
	"strings"

	"net/http"
	_ "net/http/pprof"

	"github.com/golang/glog"
	"github.com/sbezverk/tools"
	"github.com/sebasttiano/gobmp/pkg/dumper"
	"github.com/sebasttiano/gobmp/pkg/filer"
	"github.com/sebasttiano/gobmp/pkg/gobmpsrv"
	"github.com/sebasttiano/gobmp/pkg/kafka"
	"github.com/sebasttiano/gobmp/pkg/nats"
	"github.com/sebasttiano/gobmp/pkg/pub"
)

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		glog.Error("parsing config failed")
		return
	}

	// Starting performance collecting http server
	go func() {
		glog.Info(http.ListenAndServe(fmt.Sprintf(":%d", cfg.PerfPort), nil))
	}()
	// Initializing publisher
	var publisher pub.Publisher
	//var err error
	switch strings.ToLower(cfg.Dump) {
	case "file":
		publisher, err = filer.NewFiler(cfg.File)
		if err != nil {
			glog.Errorf("failed to initialize file publisher with error: %+v", err)
			os.Exit(1)
		}
		glog.V(5).Infof("file publisher has been successfully initialized.")
	case "console":
		publisher, err = dumper.NewDumper()
		if err != nil {
			glog.Errorf("failed to initialize console publisher with error: %+v", err)
			os.Exit(1)
		}
		glog.V(5).Infof("console publisher has been successfully initialized.")
	case "nats":
		publisher, err = nats.NewPublisher(cfg.NatsSrv)
		if err != nil {
			glog.Errorf("failed to initialize NATS publisher with error: %+v", err)
			os.Exit(1)
		}
		glog.V(5).Infof("NATS publisher has been successfully initialized.")
	default:

		if cfg.KafkaTopic != "" {
			publisher, err = kafka.NewKafkaSinglePublisher(cfg.KafkaUser, cfg.KafkaPass, cfg.KafkaTopic, cfg.KafkaSrv)
		} else {
			publisher, err = kafka.NewKafkaPublisher(cfg.KafkaSrv)
		}
		if err != nil {
			glog.Errorf("failed to initialize Kafka publisher with error: %+v", err)
			os.Exit(1)
		}
		glog.V(5).Infof("Kafka publisher has been successfully initialized.")
	}

	// Initializing bmp server
	interceptFlag, err := strconv.ParseBool(cfg.Intercept)
	if err != nil {
		glog.Errorf("failed to parse to bool the value of the intercept flag with error: %+v", err)
		os.Exit(1)
	}
	splitAFFlag, err := strconv.ParseBool(cfg.SplitAF)
	if err != nil {
		glog.Errorf("failed to parse to bool the value of the intercept flag with error: %+v", err)
		os.Exit(1)
	}
	bmpSrv, err := gobmpsrv.NewBMPServer(cfg.SrcPort, cfg.DstPort, interceptFlag, publisher, splitAFFlag)
	if err != nil {
		glog.Errorf("failed to setup new gobmp server with error: %+v", err)
		os.Exit(1)
	}
	// Starting Interceptor server
	bmpSrv.Start()

	stopCh := tools.SetupSignalHandler()
	<-stopCh

	bmpSrv.Stop()
	os.Exit(0)
}
