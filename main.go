package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/seanhood/go-vedirect/vedirect"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

// Config is where we keep the flag vars
type Config struct {
	device  string
	outFile string
	verbose bool
	ver     bool

	MQTT struct {
		Server  string
		Topic   string
		TLSKey  string
		TLSCert string
		TLSCA   string
	}
}

func main() {
	c := new(Config)
	flag.StringVar(&c.device, "dev", "/dev/ttyUSB0", "full path to serial device node")
	flag.StringVar(&c.MQTT.Server, "mqtt.server", "tcp://localhost:1883", "MQTT Server address")
	flag.StringVar(&c.MQTT.Topic, "mqtt.topic", "", "The MQTT Topic to publish messages to")

	flag.StringVar(&c.MQTT.TLSKey, "mqtt.tls_key", "", "MQTT TLS Private Key")
	flag.StringVar(&c.MQTT.TLSCert, "mqtt.tls_cert", "", "MQTT TLS Private Cert")
	flag.StringVar(&c.MQTT.TLSCA, "mqtt.tls_rootca", "", "MQTT TLS Root CA")

	flag.StringVar(&c.outFile, "out-file", "", "File to write json data to")
	flag.BoolVar(&c.verbose, "verbose", false, "Verbose Output")

	flag.BoolVar(&c.ver, "v", false, "Print Version")
	flag.Parse()

	if c.ver {
		fmt.Println(buildVersion(version, commit, date))
		os.Exit(0)
	}

	// Mqtt Setup
	var mqttClient mqtt.Client
	if c.MQTT.Topic != "" {
		opts := *mqtt.NewClientOptions()
		opts.SetMaxReconnectInterval(1 * time.Second)

		certpool := x509.NewCertPool()
		pemCerts, err := ioutil.ReadFile(c.MQTT.TLSCA)
		if err == nil {
			certpool.AppendCertsFromPEM(pemCerts)
		}

		if c.MQTT.TLSCert != "" && c.MQTT.TLSKey != "" {
			cert, err := tls.LoadX509KeyPair(c.MQTT.TLSCert, c.MQTT.TLSKey)
			if err != nil {
				log.Fatal(err)
			}

			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
				ClientAuth:         tls.NoClientCert,
				Certificates:       []tls.Certificate{cert},
				ClientCAs:          nil,
				RootCAs:            certpool,
			}
			opts.SetTLSConfig(tlsConfig)
		}
		opts.AddBroker(c.MQTT.Server)

		mqttClient = mqtt.NewClient(&opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			return
		}

		log.Printf("Connected to %s\n", c.MQTT.Server)

	}

	// file output setup
	var outFileHandle *os.File
	if c.outFile != "" {
		log.Printf("Saving data to: %s", c.outFile)
		var err error
		outFileHandle, err = os.OpenFile(c.outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatal(err)
		}

		defer outFileHandle.Close()

	}

	var reader io.Reader

	stat, err := os.Stat(c.device)
	if err != nil {
		log.Fatal(err)
	}

	// Should probably be in go-vedirect package
	if stat.Mode().IsRegular() {
		reader = vedirect.OpenFile(c.device)
	} else {
		reader = vedirect.OpenSerial(c.device)
	}

	s := vedirect.NewStream(reader)
	for {
		b, checksum := s.ReadBlock()
		if checksum == 0 {

			fields := b.Fields()

			fields["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)

			jsonPayload, err := json.Marshal(fields)
			if err != nil {
				log.Fatal(err)
			}

			if c.verbose {
				log.Println(string(jsonPayload))
			}

			if c.MQTT.Topic != "" {
				mqttClient.Publish(c.MQTT.Topic, 1, false, jsonPayload)
			}

			if c.outFile != "" {
				_, err := outFileHandle.Write(jsonPayload)
				if err != nil {
					log.Fatal(err)
				}
			}

		} else {
			log.Println("Bad block, skipping:", b)
		}
	}
}

func buildVersion(version, commit, date string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	return result
}
