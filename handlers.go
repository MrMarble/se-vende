package sevende

import (
	"fmt"

	"gihub.com/mrmarble/se-vende/providers"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/telebot.v3"
)

// handleStart triggers when /start is sent on private
func (t *Telegram) handleStart(ctx tb.Context) error {
	m := ctx.Message()
	if !m.Private() {
		return nil
	}

	t.send(m.Chat, "Agrega el bot al grupo que quieras como administrador y envía un enlace de Wallapop para obtener información sobre el producto")

	return nil
}

// handleQuery triggers when a query is sent
func (t *Telegram) handleQuery(ctx tb.Context) error {
	if ctx.Query().Text == "" {
		return nil
	}
	log.Info().Str("module", "telegram").Str("url", ctx.Query().Text).Msg("query received")
	results := make(tb.Results, 1)

	if url := providers.GetURL(ctx.Query().Text); url == "" {
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
			ThumbURL:    product.Images[0],
			Title:       product.Name,
			Description: product.Desc,
			Text:        product.String(),
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
	// Ignore own messages
	if ctx.Sender().ID == t.bot.Me.ID {
		return nil
	}

	m := ctx.Message()
	if m.Private() {
		t.send(m.Chat, "Este bot solo funciona en grupos")
		return nil
	}

	if url := providers.GetURL(m.Text); url != "" {
		product, err := NewProduct(url)
		if err != nil {
			log.Error().Str("module", "telegram").Err(err).Msg("error fetching product")
			ctx.Reply("Error al obtener el producto")
			return err
		}

		canDelete := true
		err = ctx.Delete()
		if err != nil {
			canDelete = false
			log.Error().Str("module", "telegram").Err(err).Msg("error deleting message")
		}

		msg := product.String()
		msg += fmt.Sprintf("🙋🏻‍♂️ <b>Contacta con el vendedor:</b> %s", getUser(ctx.Sender()))
		album := tb.Album{}
		for i, img := range product.Images {
			if i == 0 {
				album = append(album, &tb.Photo{File: tb.FromURL(img), Caption: msg})
			} else {
				album = append(album, &tb.Photo{File: tb.FromURL(img)})
			}
		}

		opts := tb.SendOptions{
			ParseMode:             tb.ModeHTML,
			DisableWebPagePreview: true,
		}
		if m.ReplyTo != nil {
			opts.ReplyTo = m.ReplyTo
		}

		if !canDelete {
			opts.ReplyTo = m
		}
		t.bot.SendAlbum(m.Chat, album, &opts)
	}
	return nil
}
