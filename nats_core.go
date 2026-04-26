package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func core() {
	slog.Info("Running Core test")
	nc, close, err := GetNatsConn()
	if err != nil {
		return
	}
	defer close()

	err = nc.Publish("first-message", []byte("Hello World!"))
	if err != nil {
		slog.Warn("Failed to send First message", "err", err)
	}

	sub, err := nc.Subscribe("request", func(m *nats.Msg) {
		slog.Info("Got message on [request]", "message", string(m.Data))
		err = m.Respond([]byte("Answer is 42"))
		if err != nil {
			slog.Warn("Failed to Respond to message", "err", err)
			return
		}

	})
	if err != nil {
		slog.Warn("Failed to Subscribe", "err", err)
		return
	}

	defer func() {
		_ = sub.Unsubscribe()
		_ = sub.Drain()
	}()

	msg, err := nc.Request("request", []byte("What is the Answer?"), 10*time.Millisecond)
	if err != nil {
		slog.Warn("Failed to send message with Subscriber", "err", err)
		return
	}

	slog.Info("Got Back Message", "message", string(msg.Data))

	_, _ = fmt.Fprintf(os.Stdout, "Press [Enter] to end\n")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')

}
