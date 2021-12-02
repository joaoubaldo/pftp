package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/DataDog/datadog-go/v5/statsd"
	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/pyama86/pftp/example/webapi"
	"github.com/pyama86/pftp/pftp"
	"github.com/sirupsen/logrus"
)

type config struct {
	Statsd statsdConfig `toml:"statsd"`
}

type statsdConfig struct {
	Host   string `toml:"host"`
	Prefix string `toml:"prefix"`
}

var confFile = "./config.toml"

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	logrus.AddHook(logrus_stack.NewHook(stackLevels, stackLevels))
}

func main() {
	ftpServer, err := pftp.NewFtpServer(confFile)
	eventC := pftp.NewEventChan(0)
	ftpServer.SetEventC(eventC)
	if err != nil {
		logrus.Fatal(err)
	}

	go func() {
		var conf config
		_, err := toml.DecodeFile(confFile, &conf)
		if err != nil {
			logrus.Errorf("Statsd goroutine failed to start: %v", err.Error())
			return
		}

		statsd, err := statsd.New(conf.Statsd.Host)
		if err != nil {
			logrus.Errorf("Statsd goroutine failed to start: %v", err.Error())
			return
		}
		defer statsd.Close()

		for {
			ev := <-eventC
			logrus.Printf("Received event: %s: %s", ev.Name(), ev.Payload())
			switch ev.Name() {
			case pftp.ClientCommandEventType:
			case pftp.ClientConnectEventType:
			case pftp.ClientDisconnectEventType:
			case pftp.DataTransferEventType:
			case pftp.ErrorEventType:
				statsd.Count("errors", 1, []string{""}, 1.0)
			}
		}
	}()

	ftpServer.Use("user", User)
	if err := ftpServer.Start(); err != nil {
		logrus.Fatal(err)
	}
}

// User function will setup Origin ftp server domain from ftp username
// If failed get domain from server, the origin will set by local (localhost:21)
func User(c *pftp.Context, param string) error {
	res, err := webapi.GetDomainFromWebAPI(confFile, param)
	if err != nil {
		logrus.Debug(fmt.Sprintf("cannot get origin host from webapi server:%v", err))
		c.RemoteAddr = ""
	} else {
		c.RemoteAddr = *res
	}

	return nil
}
