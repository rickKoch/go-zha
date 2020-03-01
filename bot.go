package zha

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Adapter interface
type Adapter interface {
	Register(*Brain)
	Send(text, channelID string) error
	Close() error
}

// Bot struct
type Bot struct {
	Context context.Context
	Name    string
	Adapter Adapter
	Memory  Memory
	Brain   *Brain
	Logger  *zap.Logger

	initErr error
}

// NewBot generates new bot
func NewBot(name string, opts ...Option) *Bot {

	b := &Bot{
		Context: context.Background(),
		Logger:  NewLogger(),
		Name:    name,
	}

	timeout := 10 * time.Second
	b.Brain = NewBrain(b.Logger, timeout)
	b.Logger.Info("Init bot", zap.String("name", name))

	for _, opt := range opts {
		err := opt(b)
		if err != nil && b.initErr == nil {
			b.initErr = err
		}
	}

	if b.Memory == nil {
		b.Memory = NewInMemory()
	}

	return b

}

// Run starts the bot
func (b *Bot) Run() error {
	if b.initErr != nil {
		return errors.Wrap(b.initErr, "failed to init bot")
	}

	b.Adapter.Register(b.Brain)
	b.Brain.Emit(InitEvent{})

	b.Logger.Info("Bot initialized and ready to operate", zap.String("name", b.Name))

	b.Brain.Process(b.Context)

	err := b.Adapter.Close()
	b.Logger.Info("Bot is shuthig down", zap.String("name", b.Name))
	if err != nil {
		b.Logger.Info("Error while closing adapter", zap.Error(err))
	}

	return nil
}

// Respond gets the message and register the handlers
func (b *Bot) Respond(msg string, fun func(Message) error) {
	expr := "^" + msg + "$"
	if expr == "" {
		return
	}

	if expr[0] == '^' {
		if !strings.HasPrefix(expr, "^(?i)") {
			expr = "^(?i)" + expr[1:]
		}
	} else {
		if !strings.HasPrefix(expr, "(?i)") {
			expr = "(?i)" + expr
		}
	}

	regex, err := regexp.Compile(expr)
	if err != nil {
		b.Logger.Error("Failed to add Response handler", zap.Error(err))
		return
	}

	b.Brain.RegisterHandler(func(ctx context.Context, evt ReciveMessageEvent) error {
		matches := regex.FindStringSubmatch(evt.Text)
		if len(matches) == 0 {
			return nil
		}

		return fun(Message{
			Context:  ctx,
			Text:     evt.Text,
			ChannelD: evt.ChannelD,
			Matches:  matches[1:],
			adapter:  b.Adapter,
		})
	})
}
