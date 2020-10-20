# go-vedirect-publisher
A utility for publishing [VE.Direct](https://www.victronenergy.com/live/vedirect_protocol:faq) frames over MQTT. With support for embedded devices (Teltonika RUT)

This program is for reading data from a Victron device using the VE.Direct protocol.
It's built for the Teltonika RUT955 Router, however likely can be compiled for any device which Go can be compiled for. The RUT955 has a MIPS processor which runs RutOS, a fork of OpenWrt.

The two devices are connected together using the [VE.Direct to USB interface](https://www.victronenergy.com/accessories/ve-direct-to-usb-interface)

go-vedirect-publisher has currently been tested with Mosquitto MQTT Server and AWS IoT Core

### Example payload

```
{"CS":"0","ERR":"0","FW":"150","H19":"325","H20":"0","H21":"0","H22":"0","H23":"0","HSDS":"213","I":"-230","IL":"200","LOAD":"ON","MPPT":"0","OR":"0x00000001","PID":"0xA053","PPV":"0","SER#":"XXXXXXXXXXX","V":"13210","VPV":"10","timestamp":"1605565193"}
```

## Features
- Send payload to an MQTT Server
- Save payload to file

## Usage
```
Usage of ./bin/go-vedirect:
  -dev string
		full path to serial device node (default "/dev/ttyUSB0")
  -mqtt.server string
		MQTT Server address (default "tcp://localhost:1883")
  -mqtt.tls_cert string
		MQTT TLS Private Cert
  -mqtt.tls_key string
		MQTT TLS Private Key
  -mqtt.tls_rootca string
		MQTT TLS Root CA
  -mqtt.topic string
		The MQTT Topic to publish messages to
  -out-file string
		File to write json data to
  -v	Print Version
  -verbose
		Verbose Output
```

## Roadmap
- Tests
- Configurable sample rate, vedirect updates every second
- Send payload to HTTP endpoint
- Rewrite in better language for embedded systems (C++)
