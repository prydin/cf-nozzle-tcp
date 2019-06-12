package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kelseyhightower/envconfig"
	"github.com/prydin/cf-nozzle-tcp/nozzle"
)

func main() {
	nozzleConfig := &nozzle.NozzleConfig{}
	err := envconfig.Process("nozzle", nozzleConfig)
	if err != nil {
		log.Fatalf("Error processing environment: %s", err)
	}

	client, err := nozzle.NewClient(nozzleConfig)
	if err != nil {
		log.Fatalf("Error obtaining security token: %s", err)
	}

	cnsmr := consumer.New(client.Endpoint.DopplerEndpoint, &tls.Config{InsecureSkipVerify: nozzleConfig.SkipSSL}, nil)
	cnsmr.SetDebugPrinter(ConsoleDebugPrinter{})

	var (
		msgChan   <-chan *events.Envelope
		errorChan <-chan error
	)

	token, err := client.GetToken()
	if err != nil {
		log.Fatalf("Error extracting security token: %s", err)
	}

	msgChan, errorChan = cnsmr.FilteredFirehose(nozzleConfig.FirehoseSubscriptionID, token, consumer.LogMessages)

	go func() {
		for err := range errorChan {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
	}()

	sender := nozzle.NewTCPSender(nozzleConfig.Target)

	for msg := range msgChan {
		txt := fmt.Sprintf("origin: %s, job: %s, index: %s, ip: %s, message: %s\n",
			*msg.Origin, *msg.Job, *msg.Index, *msg.Ip, string(msg.LogMessage.Message))
		if nozzleConfig.Debug {
			fmt.Print(txt)
		}
		sender.Send(txt)
	}
}

type ConsoleDebugPrinter struct{}

func (c ConsoleDebugPrinter) Print(title, dump string) {
	println(title)
	println(dump)
}
