package main

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/kochevRisto/go-zha"
	"gitlab.com/kochevRisto/go-zha/redis-memory"
	"gitlab.com/kochevRisto/go-zha/slack"
)

type ExampleBot struct {
	*zha.Bot
}

func main() {

	adapter := slack.NewAdapter("xoxb-953511947447-940297123539-R4O3Cxt6W2Od0UuEQrOTXqwt")
	bot := &ExampleBot{
		Bot: zha.NewBot(
			"test",
			// file.MemoryOption("./test.json")),
			redis.Memory("localhost:6379", redis.WithKey("risto-bot")),
		),
	}

	adapter(bot.Bot)

	bot.Respond("remember (.+) is (.+)", bot.Remember)
	bot.Respond(`(what is) ([^?]+)\s*\??(.*)`, bot.WhatIs)
	bot.Respond(`forget (.+)`, bot.Forget)
	bot.Respond(`(.*)what do you remember\??(.*)`, bot.WhatDoYouRemember)

	bot.Run()

}

// Remember a value for a given key.
//   command: bot remember <key> is <value>
func (b *ExampleBot) Remember(msg zha.Message) error {
	key, value := msg.Matches[0], msg.Matches[1]
	key = strings.TrimSpace(key)
	msg.Respond("\nOK, I'll remember %s is %s\n", key, value)
	return b.Memory.Set(key, value)
}

// WhatIs test
func (b *ExampleBot) WhatIs(msg zha.Message) error {
	key := strings.TrimSpace(msg.Matches[1])
	value, ok, err := b.Memory.Get(key)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve key %q from brain", key)
	}

	if ok {
		msg.Respond("\n%s is %s\n", key, value)
	} else {
		msg.Respond("\nI do not remember %q\n", key)
	}

	return nil
}

// Forget test
func (b *ExampleBot) Forget(msg zha.Message) error {
	key := strings.TrimSpace(msg.Matches[0])
	value, _, _ := b.Memory.Get(key)
	ok, err := b.Memory.Delete(key)
	if err != nil {
		return errors.Wrapf(err, "failed to delete key %q from brain", key)
	}

	if !ok {
		msg.Respond("\n I do not remember %q\n", key)
	} else {
		msg.Respond("\n I've forgotten %s is %s.\n", key, value)
	}

	return nil
}

// WhatDoYouRemember test
func (b *ExampleBot) WhatDoYouRemember(msg zha.Message) error {
	data, err := b.Memory.Memories()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve all memories from brain")
	}

	switch len(data) {
	case 0:
		msg.Respond("\nI do not remember anything\n")
		return nil
	case 1:
		msg.Respond("\nI have only a single memory:\n")
	default:
		msg.Respond("\nI have %d memories:\n", len(data))
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		value := data[key]
		msg.Respond("\n%s is %s\n", key, value)
	}

	return nil
}
