package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

const (
	MainURL = "nats://gnats.laserath.com:4224"
	ReplURL = "nats://gnats-repl.laserath.com:4225"
)

func GetNatsConn() (*nats.Conn, func(), error) {
	connString := fmt.Sprintf("%s,%s", MainURL, ReplURL)
	nc, err := nats.Connect(connString)
	if err != nil {
		slog.Info("Failed to connect to Gnats Cluster", "err", err)
		return nil, func() {}, err
	}
	close := func() {
		_ = nc.Drain()
		nc.Close()
	}

	return nc, close, nil

}

func main() {
	var proc string
	flag.StringVar(&proc, "proc", "core", "type of test to run [core/jetstream]")

	flag.Parse()

	switch proc {
	case "core":
		core()
	case "jetstream":
		Jetstream()
	default:
		slog.Error("Need to specify proc as either core or jetstream")
	}
}
