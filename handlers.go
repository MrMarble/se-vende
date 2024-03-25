package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/telebot.v3"
)

// handleStart triggers when /start is sent on private
func (t *Telegram) handleStart(ctx tb.Context) error {
	m := ctx.Message()
	if !m.Private() {
		return nil
	}

	t.send(m.Chat, "enyaki gay")

	return nil
}

// handleQuery triggers when a query is sent
func (t *Telegram) handleQuery(ctx tb.Context) error {
	if ctx.Query().Text == "" {
		return nil
	}
	log.Info().Str("module", "telegram").Str("url", ctx.Query().Text).Msg("query received")
	results := make(tb.Results, 1)

	if url := hasValidURL(ctx.Query().Text); url == "" {
		results[0] = &tb.ArticleResult{
			Title:       "Error",
			Text:        "URL inválida",
			Description: "URL inválida",
		}
	} else {
		product, err := NewProduct(url)
		if err != nil {
			log.Error().Str("module", "telegram").Err(err).Msg("error fetching product")
			results[0] = &tb.ArticleResult{
				Title: "Error",
				Text:  "Error al obtener el producto",
			}
			return ctx.Answer(&tb.QueryResponse{
				Results:   results,
				CacheTime: 10,
			})
		}
		r := t.bot.NewMarkup()
		r.Inline(r.Row(r.URL("Ver en Wallapop", url)))
		log.Info().Str("module", "telegram").Interface("product", product).Msg("product fetched")
		result := &tb.ArticleResult{
			ThumbURL:    product.Image,
			Title:       product.Name,
			Description: product.Desc,
			Text:        product.String(),
			URL:         product.Image,
			ResultBase: tb.ResultBase{
				ParseMode:   tb.ModeHTML,
				ReplyMarkup: r,
			},
		}
		results[0] = result
	}

	return ctx.Answer(&tb.QueryResponse{
		Results:   results,
		CacheTime: 10,
	})
}

// handleText triggers when a text message is sent
func (t *Telegram) handleText(ctx tb.Context) error {
	m := ctx.Message()
	if m.Private() {
		t.send(m.Chat, "Este bot solo funciona en grupos")
		return nil
	}

	if url := hasValidURL(m.Text); url != "" {
		product, err := NewProduct(url)
		if err != nil {
			log.Error().Str("module", "telegram").Err(err).Msg("error fetching product")
			t.send(m.Chat, "Error al obtener el producto")
			return err
		}
		ctx.Delete()

		msg := product.String()
		msg += fmt.Sprintf("🙋🏻‍♂️ <b>Contacta con el vendedor:</b> %s", getUser(ctx.Sender()))
		t.send(m.Chat, msg, tb.ModeHTML)
	} else {
		t.send(m.Chat, "URL inválida")
	}
	return nil
}
