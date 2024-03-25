package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

// Response represents a wallapop item
type Response struct {
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
						Amount int `json:"amount"`
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

// NewItem returns a new item
func newItem(url string) (*Response, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Error().Str("module", "wallapop").Err(err).Msg("error fetching url")
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error().Str("module", "wallapop").Msgf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("error loading the document")
		return nil, err
	}

	raw := doc.Find("script#__NEXT_DATA__").Text()
	if raw == "" {
		log.Error().Str("module", "wallapop").Msg("error fetching product")
		return nil, fmt.Errorf("error fetching product")
	}

	// decode the json
	var data Response
	err = json.Unmarshal([]byte(raw), &data)
	if err != nil {
		log.Error().Str("module", "wallapop").Err(err).Msg("error decoding json")
		return nil, err
	}

	return &data, nil
}
