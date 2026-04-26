package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func Jetstream() {
	slog.Info("Running Jestream test")

	nc, close, err := GetNatsConn()
	if err != nil {
		return
	}
	defer close()

	js, err := jetstream.New(nc)
	if err != nil {
		slog.Warn("Failed to initalize jetstream", "err", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "mystream",
		Subjects:  []string{"mystream.stream_sub.>"},
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		slog.Warn("Failed to setup Stream", "err", err)
		return
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		slog.Warn("Failed to setup Consumer", "err", err)
		return
	}

	consumeHandle, _ := consumer.Consume(func(msg jetstream.Msg) {
		slog.Info("Got Jetstream Message", "msg", string(msg.Data()))
		err := msg.Ack()
		if err != nil {
			slog.Warn("Failed to send Ack", "err", err)
		}

	})

	defer consumeHandle.Stop()

	for i := range 10 {
		data := []byte(fmt.Sprintf("Message Numero %d", i))
		ack, err := js.PublishMsg(ctx, &nats.Msg{
			Data:    data,
			Subject: "mystream.stream_sub.numeros",
		})
		if err != nil {
			slog.Warn("Failed to send message", "numero", i)
		}

		slog.Info("Got Ack", "ack", ack.Value, "numero", i)
	}
}
