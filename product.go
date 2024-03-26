package sevende

import (
	"fmt"
	"strings"

	"gihub.com/mrmarble/se-vende/providers"
)

type Product struct {
	Provider  string
	URL       string
	Name      string
	Desc      string
	Condition string
	Images    []string
	Tags      []string
	Price     float32
	Location  string
	Shipping  bool
}

func NewProduct(url string) (*Product, error) {
	item, err := providers.NewItem(url)
	if err != nil {
		return nil, err
	}

	return &Product{
		Provider:  item.Provider,
		URL:       item.URL,
		Name:      item.Title,
		Desc:      item.Description,
		Price:     item.Price,
		Location:  item.Location,
		Shipping:  item.Shipping,
		Condition: item.Status,
		Images:    item.Images,
		Tags:      item.Tags,
	}, nil
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
		output += fmt.Sprintf("<b>💶 Precio:</b> %.2f€\n", p.Price)
	}

	tags := strings.Join(p.Tags, " #")
	if len(tags) > 0 {
		output += fmt.Sprintf("<b>🏷 Etiquetas:</b> #%s\n", tags)
	}

	output += fmt.Sprintf("\n<b>🔗 Enlace a <a href='%s'>%s</a></b>\n", p.URL, p.Provider)
	return output
}
