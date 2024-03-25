package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Product struct {
	URL       string
	Name      string
	Desc      string
	Condition string
	Images    []string
	Tags      []string
	Price     int
	Location  string
	Shipping  bool
}

func NewProduct(url string) (*Product, error) {
	item, err := newItem(url)
	if err != nil {
		return nil, err
	}

	product := &Product{
		URL:       url,
		Name:      item.Props.PageProps.Item.Title.Original,
		Desc:      item.Props.PageProps.Item.Description.Original,
		Condition: item.Props.PageProps.Item.Condition.Text,
		Price:     item.Props.PageProps.Item.Price.Cash.Amount,
		Location:  fmt.Sprintf("%s (%s)", item.Props.PageProps.Item.Location.City, item.Props.PageProps.Item.Location.PostalCode),
		Shipping:  item.Props.PageProps.Item.Shipping.ItemShippable && item.Props.PageProps.Item.Shipping.ShippingAllowedByUser,
		Tags:      item.Props.PageProps.Item.Hashtags,
	}

	for _, img := range item.Props.PageProps.Item.Images {
		product.Images = append(product.Images, img.URLs.Small)
	}

	return product, nil
}

func hasValidURL(text string) string {
	regex := regexp.MustCompile(`(https:\/\/([a-z]+\.)?wallapop\.com\/item\/.*)`)
	return regex.FindString(text)
}

func (p *Product) String() string {
	output := fmt.Sprintf("<b>🔎 Producto:</b> %s\n\n", p.Name)
	if p.Condition != "" {
		output += fmt.Sprintf("<b>🔧 Estado:</b> %s\n", p.Condition)
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

	tags := strings.Join(p.Tags, " #")
	if len(tags) > 0 {
		output += fmt.Sprintf("<b>🏷 Etiquetas:</b> #%s\n", tags)
	}

	output += fmt.Sprintf("\n<b>🔗 Enlace a <a href='%s'>Wallapop</a></b>\n", p.URL)
	return output
}
