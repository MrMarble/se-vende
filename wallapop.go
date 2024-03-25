package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

type Product struct {
	URL             string
	Image           string
	Characteristics string
	Tags            []string
	Name            string
	Desc            string
	Price           uint32
	Location        string
	Shipping        bool
}

func NewProduct(url string) (*Product, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("error fetching url")
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error().Msgf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("error loading the document")
		return nil, err
	}

	//fmt.Println(doc.Find("script#__NEXT_DATA__").Html())

	article := doc.Find("[class^='item-detail_ItemDetail__card']")

	regex := regexp.MustCompile(`\d+`)
	price, err := strconv.Atoi(string(regex.Find([]byte(article.Find("[class^='item-detail-price_ItemDetailPrice']").Text()))))
	if err != nil {
		log.Error().Err(err).Msg("error parsing price")
		return nil, err
	}

	img := doc.Find("meta[name='twitter:image']").AttrOr("content", "")
	name := article.Find("[class^='item-detail_ItemDetail__title']").Text()
	desc := article.Find("[class^='item-detail_ItemDetail__description']").Text()
	location := article.Find("[class^='d-flex item-detail-location_ItemDetailLocation']").Text()
	shipping := article.Find("[badge-type='shippingAvailable']").Length() > 0
	characteristics := article.Find("[class^='item-detail-additional-specifications_ItemDetailAdditionalSpecifications']").Text()
	tags := []string{}
	article.Find("[class^='item-detail_ItemDetail__separator__']:first-child div span").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})

	return &Product{
		URL:             url,
		Name:            name,
		Desc:            desc,
		Image:           img,
		Characteristics: characteristics,
		Price:           uint32(price),
		Tags:            tags,
		Location:        location,
		Shipping:        shipping,
	}, nil
}

func hasValidURL(text string) string {
	regex := regexp.MustCompile(`(https:\/\/([a-z]+\.)?wallapop\.com\/item\/.*)`)
	return regex.FindString(text)
}

func (p *Product) String() string {
	output := fmt.Sprintf("<b><a href='%s'>🔎</a> Producto:</b> %s\n\n", p.Image, p.Name)
	if p.Characteristics != "" {
		output += fmt.Sprintf("<b>🔧 Estado:</b> %s\n", p.Characteristics)
	}
	if p.Desc != "" {
		output += fmt.Sprintf("<b>🏺 Descripción:</b> %s\n\n", p.Desc)
	}
	if p.Location != "" {
		output += fmt.Sprintf("<b>🌍 Localidad:</b> %s\n", p.Location)
	}

	ship := "❌"
	if p.Shipping {
		ship = "✅"
	}
	output += fmt.Sprintf("<b>📦 Envío:</b> %s\n", ship)

	if p.Price != 0 {
		output += fmt.Sprintf("<b>💶 Precio:</b> %d€\n", p.Price)
	}

	tags := strings.Join(p.Tags, " ")
	if len(tags) > 0 {
		output += fmt.Sprintf("<b>🏷 Etiquetas:</b> %s\n", tags)
	}

	output += fmt.Sprintf("\n<b>🔗 Enlace a <a href='%s'>Wallapop</a></b>\n", p.URL)
	return output
}
