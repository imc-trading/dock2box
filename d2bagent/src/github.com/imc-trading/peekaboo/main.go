package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Unknwon/macaron"
	flags "github.com/jessevdk/go-flags"
	"github.com/mickep76/hwinfo"
)

func main() {
	// Set log options.
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	// Options.
	var opts struct {
		Verbose      bool    `short:"v" long:"verbose" description:"Verbose"`
		Version      bool    `long:"version" description:"Version"`
		BindAddr     string  `short:"b" long:"bind-addr" description:"Bind to address" default:"0.0.0.0"`
		Port         int     `short:"p" long:"port" description:"Port" default:"5050"`
		StaticDir    string  `short:"s" long:"static-dir" description:"Static content" default:"static"`
		TemplateDir  string  `short:"t" long:"template-dir" description:"Templates" default:"templates"`
		KafkaEnabled bool    `short:"K" long:"kafka" description:"Enable Kafka message bus"`
		KafkaTopic   string  `long:"kafka-topic" description:"Kafka topic" default:"peekaboo"`
		KafkaPeers   *string `long:"kafka-peers" description:"Comma-delimited list of Kafka brokers"`
		KafkaCert    *string `long:"kafka-cert" description:"Certificate file for client authentication"`
		KafkaKey     *string `long:"kafka-key" description:"Key file for client client authentication"`
		KafkaCA      *string `long:"kafka-ca" description:"CA file for TLS client authentication"`
		KafkaVerify  bool    `long:"kafka-verify" description:"Verify SSL certificate"`
	}

	// Parse options.
	if _, err := flags.Parse(&opts); err != nil {
		ferr := err.(*flags.Error)
		if ferr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(err.Error())
		}
	}

	// Print version.
	if opts.Version {
		fmt.Printf("peekaboo %s\n", Version)
		os.Exit(0)
	}

	// Set verbose.
	if opts.Verbose {
		log.SetLevel(log.InfoLevel)
	}

	// Check root.
	if runtime.GOOS != "darwin" && os.Getuid() != 0 {
		log.Fatal("This application requires root privileges to run.")
	}

	// Get hardware info.
	info := hwinfo.NewHWInfo()
	if err := info.GetTTL(); err != nil {
		log.Fatal(err.Error())
	}

	// Produce message to Kafka bus.
	log.Infof("Produce startup event to Kafka bus with topic: %s", opts.KafkaTopic)
	if opts.KafkaEnabled {
		if opts.KafkaPeers == nil {
			log.Fatal("You need to specify Kafka Peers")
		}

		user, err := user.Current()
		if err != nil {
			log.Fatal(err.Error())
		}

		event := &Event{
			Name:      "Peekaboo startup",
			EventType: STARTED,
			Created:   time.Now().Format("20060102T150405ZB"),
			CreatedBy: CreatedBy{
				User:    user.Username,
				Service: "peekaboo",
				Host:    info.Hostname,
			},
			Descr: "Peekaboo startup event",
			Data:  info,
		}

		producer := newProducer(strings.Split(*opts.KafkaPeers, ","), opts.KafkaCert, opts.KafkaKey, opts.KafkaCA, opts.KafkaVerify)
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: opts.KafkaTopic,
			Value: event,
		})
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Infof("Kafka partition: %v, offset: %v", partition, offset)
	}

	log.Infof("Using static dir: %s", opts.StaticDir)
	log.Infof("Using template dir: %s", opts.TemplateDir)

	m := macaron.Classic()
	m.Use(macaron.Static(opts.StaticDir))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  opts.TemplateDir,
		IndentJSON: true,
	}))

	routes(m, info)
	m.Run(opts.BindAddr, opts.Port)
}
