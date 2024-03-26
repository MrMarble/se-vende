package sevende

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/telebot.v3"
)

// Telegram represents the telegram bot
type Telegram struct {
	bot                *tb.Bot
	handlersRegistered bool
}

// NewBot returns a Telegram bot
func NewBot(token string) (*Telegram, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, _ tb.Context) {
			log.Error().Str("module", "telegram").Err(err).Msg("telebot internal error")
		},
	})

	if err != nil {
		return nil, err
	}

	log.Info().Str("module", "telegram").Int64("id", bot.Me.ID).Str("name", bot.Me.FirstName).Str("username", bot.Me.Username).Msg("connected to telegram")

	return &Telegram{bot: bot}, nil
}

// Start starts polling for telegram updates
func (t *Telegram) Start() {
	t.registerHandlers()

	log.Info().Str("module", "telegram").Msg("start polling")
	t.bot.Start()
}

// RegisterHandlers registers all the handlers
func (t *Telegram) registerHandlers() {
	if t.handlersRegistered {
		return
	}

	log.Info().Str("module", "telegram").Msg("registering handlers")

	t.bot.Handle("/start", t.handleStart)
	t.bot.Handle(tb.OnText, t.handleText)
	t.bot.Handle(tb.OnQuery, t.handleQuery)
	t.handlersRegistered = true
}

// send sends a message with error logging and retries
func (t *Telegram) send(to tb.Recipient, what interface{}, options ...interface{}) *tb.Message {
	hasParseMode := false
	for _, opt := range options {
		if _, hasParseMode = opt.(tb.ParseMode); hasParseMode {
			break
		}
	}

	if !hasParseMode {
		options = append(options, tb.ModeHTML)
	}

	try := 1
	for {
		msg, err := t.bot.Send(to, what, options...)

		if err == nil {
			return msg
		}

		if try > 5 {
			log.Error().Str("module", "telegram").Err(err).Msg("send aborted, retry limit exceeded")
			return nil
		}

		backoff := time.Second * 5 * time.Duration(try)
		log.Warn().Str("module", "telegram").Err(err).Str("sleep", backoff.String()).Msg("send failed, sleeping and retrying")
		time.Sleep(backoff)
		try++
	}
}

func getUser(u *tb.User) string {
	output := u.FirstName

	if u.LastName != "" {
		output += " " + u.LastName
	}

	if u.Username != "" {
		output += " (@" + u.Username + ")"
	} else {
		output += fmt.Sprintf(" (tg://user?id=%d)", u.ID)
	}

	return output
}
