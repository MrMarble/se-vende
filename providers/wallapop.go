package providers

import (
	"fmt"
)

// wallapop represents a wallapop item
type wallapop struct {
	Props struct {
		PageProps struct {
			Item struct {
				Title struct {
					Original string `json:"original"`
				} `json:"title"`
				Description struct {
					Original string `json:"original"`
				} `json:"description"`
				Characteristics string `json:"characteristics"`
				ShareURL        string `json:"shareUrl"`
				Images          []struct {
					URLs struct {
						Small string `json:"small"`
					} `json:"urls"`
				} `json:"images"`
				Price struct {
					Cash struct {
						Amount float32 `json:"amount"`
					} `json:"cash"`
				} `json:"price"`
				Location struct {
					City        string `json:"city"`
					CountryCode string `json:"countryCode"`
					PostalCode  string `json:"postalCode"`
				} `json:"location"`
				Shipping struct {
					ItemShippable         bool `json:"isItemShippable"`
					ShippingAllowedByUser bool `json:"isShippingAllowedByUser"`
				} `json:"shipping"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				Hashtags []string `json:"hashtags"`
			} `json:"item"`
		} `json:"pageProps"`
	} `json:"props"`
}

func fromWallapop(url string) (*Item, error) {
	resp, err := fetchItem[wallapop](url)
	if err != nil {
		return nil, err
	}

	item := &Item{
		Provider:    "wallapop",
		Title:       resp.Props.PageProps.Item.Title.Original,
		Description: resp.Props.PageProps.Item.Description.Original,
		Price:       resp.Props.PageProps.Item.Price.Cash.Amount,
		URL:         resp.Props.PageProps.Item.ShareURL,
		Location:    fmt.Sprintf("%s", resp.Props.PageProps.Item.Location.City),
		Status:      resp.Props.PageProps.Item.Condition.Text,
		Shipping:    resp.Props.PageProps.Item.Shipping.ItemShippable && resp.Props.PageProps.Item.Shipping.ShippingAllowedByUser,
		Tags:        resp.Props.PageProps.Item.Hashtags,
	}

	if resp.Props.PageProps.Item.Location.PostalCode != "" {
		item.Location = fmt.Sprintf("%s (%s)", item.Location, resp.Props.PageProps.Item.Location.PostalCode)
	}

	for _, img := range resp.Props.PageProps.Item.Images {
		item.Images = append(item.Images, img.URLs.Small)
	}

	return item, nil
}
