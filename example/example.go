package main

import (
	"fmt"
	"time"

	slogchannel "github.com/samber/slog-channel"

	"log/slog"
)

func main() {
	ch := make(chan *slog.Record, 100)
	defer close(ch)

	logger := slog.New(slogchannel.Option{Level: slog.LevelDebug, Channel: ch}.NewChannelHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("error", fmt.Errorf("an error")).
		Error("a message", slog.Int("count", 1))

	record := <-ch
	fmt.Println(record)
}
