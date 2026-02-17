package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	gray   = "\033[90m"
	blue   = "\033[34m"
)

func levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return red
	case level >= slog.LevelWarn:
		return yellow
	case level >= slog.LevelInfo:
		return green
	default:
		return reset
	}
}

type ColorHandler struct{}

func (c *ColorHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

func (c *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	color := levelColor(r.Level)

	// timestamp
	ts := r.Time.Format(time.RFC3339)

	// format attributes as key=value
	attrStr := ""
	r.Attrs(func(a slog.Attr) bool {
		attrStr += fmt.Sprintf("%s=%v ", a.Key, a.Value.Any())
		return true
	})

	// print fully formatted colored line
	fmt.Printf("%s%s [%s] %s%s%s\n",
		blue, ts, r.Level, color, r.Message, reset,
	)

	if attrStr != "" {
		fmt.Printf("    %s%s%s\n", gray, attrStr, reset)
	}

	return nil
}

func (c *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c
}

func (c *ColorHandler) WithGroup(name string) slog.Handler {
	return c
}

type TarotCard struct {
	Name    string
	Meaning string
}

func randomCard() TarotCard {
	return tarotDeck[rand.Intn(len(tarotDeck))]
}

var tarotDeck = []TarotCard{
	{"The Fool", "New beginnings, spontaneity, trust in the journey"},
	{"The Magician", "Manifestation, skill, taking action"},
	{"The High Priestess", "Intuition, inner knowledge, mystery"},
	{"The Empress", "Abundance, creativity, nurturing energy"},
	{"The Emperor", "Authority, structure, leadership"},
	{"The Lovers", "Relationships, choices, alignment of values"},
	{"The Chariot", "Determination, willpower, victory"},
	{"Strength", "Inner strength, patience, quiet confidence"},
	{"The Hermit", "Reflection, solitude, inner guidance"},
	{"Wheel of Fortune", "Change, cycles, fate"},
}

func main() {
	logger := slog.New(&ColorHandler{})
	slog.SetDefault(logger)

	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI("8557747783:AAHvDHWT6tzBh62jzw4QXO12YIr2mXWjlyE")
	if err != nil {
		logger.Error("failed to create bot", "error", err)
		os.Exit(1)
	}
	bot.Debug = false
	logger.Info("bot authorized", "username", bot.Self.UserName, "id", bot.Self.ID)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := update.Message
		chatID := msg.Chat.ID

		logger.Info(
			"incoming message",
			"chat_id", chatID,
			"user", msg.From.UserName,
			"text", msg.Text,
		)

		switch msg.Text {
		case "/start":
			logger.Info("handling /start command", "chat_id", chatID)
			reply := tgbotapi.NewMessage(
				chatID,
				"🔮 Welcome.\n\nUse /reading to draw a tarot card.",
			)
			if _, err := bot.Send(reply); err != nil {
				logger.Error("failed to send /start reply", "error", err)
			}
		case "/reading":
			card := randomCard()

			logger.Info(
				"tarot card drawn",
				"chat_id", chatID,
				"card", card.Name,
				"time", time.Now(),
			)

			text := "🃏 *Your Tarot Card*\n\n" +
				"*" + card.Name + "*\n" +
				card.Meaning

			reply := tgbotapi.NewMessage(chatID, text)
			reply.ParseMode = "Markdown"

			if _, err := bot.Send(reply); err != nil {
				logger.Error(
					"failed to send tarot reading",
					"chat_id", chatID,
					"card", card.Name,
					"error", err,
				)
			}
		default:
			logger.Warn(
				"unknown command",
				"chat_id", chatID,
				"text", msg.Text,
			)

			reply := tgbotapi.NewMessage(
				chatID,
				"Unknown command. Use /reading.",
			)
			if _, err := bot.Send(reply); err != nil {
				logger.Error("failed to send fallback reply", "error", err)
			}
		}
	}

	logger.Info("bot stopped", "context", ctx)
}
