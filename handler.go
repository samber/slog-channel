package slogchannel

import (
	"context"

	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// Channel
	Channel chan *slog.Record

	// optional: customize record builder
	Converter Converter

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func (o Option) NewChannelHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	return &ChannelHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*ChannelHandler)(nil)

type ChannelHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *ChannelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *ChannelHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	output := converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)
	h.option.Channel <- output

	return nil
}

func (h *ChannelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ChannelHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *ChannelHandler) WithGroup(name string) slog.Handler {
	return &ChannelHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
