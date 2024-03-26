package providers

import (
	"fmt"
	"strconv"
)

// vinted represents a wallapop item
type vinted struct {
	Props struct {
		PageProps struct {
			Item struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				ShareURL    string `json:"url"`
				Condition   string `json:"status"`
				Country     string `json:"country"`
				City        string `json:"city"`
				Size        string `json:"size"`
				Images      []struct {
					URL string `json:"url"`
				} `json:"photos"`
				Price struct {
					Amount string `json:"amount"`
				} `json:"total_item_price"`
			} `json:"itemDto"`
		} `json:"pageProps"`
	} `json:"props"`
}

// NewItem returns a new item
func fromVinted(url string) (*Item, error) {
	resp, err := fetchItem[vinted](url)
	if err != nil {
		return nil, err
	}

	item := &Item{
		Provider:    "vinted",
		Title:       resp.Props.PageProps.Item.Title,
		Description: resp.Props.PageProps.Item.Description,
		URL:         resp.Props.PageProps.Item.ShareURL,
		Location:    fmt.Sprintf("%s", resp.Props.PageProps.Item.Country),
		Status:      resp.Props.PageProps.Item.Condition,
		Shipping:    true,
	}

	if resp.Props.PageProps.Item.Price.Amount != "" {
		price, err := strconv.ParseFloat(resp.Props.PageProps.Item.Price.Amount, 32)
		if err != nil {
			return nil, err
		}
		item.Price = float32(price)
	}

	if resp.Props.PageProps.Item.City != "" {
		item.Location = fmt.Sprintf("%s (%s)", item.Location, resp.Props.PageProps.Item.City)
	}

	for _, img := range resp.Props.PageProps.Item.Images {
		item.Images = append(item.Images, img.URL)
	}

	if resp.Props.PageProps.Item.Size != "" {
		item.Description = fmt.Sprintf("%s\nTamaño: %s", item.Description, resp.Props.PageProps.Item.Size)
	}

	return item, nil
}
