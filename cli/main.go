package main

import (
	"net"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/nickw444/miio-go"
	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/protocol"
	"github.com/nickw444/miio-go/protocol/tokens"
	"github.com/sirupsen/logrus"
)

var sharedClient *miio.Client

func createClient(local bool) (*miio.Client, error) {
	addr := net.IPv4bcast
	if local {
		addr = net.IPv4(127, 0, 0, 1)
	}

	tokenStore, err := tokens.FromFile("tokens.txt")
	if err != nil {
		panic(err)
	}

	proto, err := protocol.NewProtocol(protocol.ProtocolConfig{
		BroadcastIP: addr,
		TokenStore:  tokenStore,
	})
	if err != nil {
		return nil, err
	}

	return miio.NewClientWithProtocol(proto)
}

func main() {
	app := kingpin.New("miio-go CLI", "CLI application to manually test miio-go functionality")
	local := app.Flag("local", "Send broadcast to 127.0.0.1 instead of 255.255.255.255 (For use with locally hosted simulator)").Bool()
	logLevel := app.Flag("log-level", "Set MiiO to a specific log level").Default("warn").Enum("debug", "warn", "info", "error")

	installControl(app)
	installDiscovery(app)

	app.Action(func(ctx *kingpin.ParseContext) error {
		level, _ := logrus.ParseLevel(*logLevel)
		l := logrus.New()
		l.SetLevel(level)
		common.SetLogger(l)

		var err error
		sharedClient, err = createClient(*local)
		return err
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
