package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

type Item struct {
	Provider    string
	Title       string
	Description string
	Price       float32
	URL         string
	Location    string
	Status      string
	Shipping    bool
	Images      []string
	Tags        []string
}

func GetURL(text string) string {
	regex := regexp.MustCompile(`((https?://)?\w*.?(wallapop|vinted)\.[comes]+\/items?\/.*)`)
	return regex.FindString(text)
}

func getProvider(url string) string {
	if url == "" {
		return ""
	}
	if regexp.MustCompile(`wallapop`).MatchString(url) {
		return "wallapop"
	} else if regexp.MustCompile(`vinted`).MatchString(url) {
		return "vinted"
	}
	return ""
}

func NewItem(url string) (*Item, error) {
	switch getProvider(url) {
	case "wallapop":
		return fromWallapop(url)
	case "vinted":
		return fromVinted(url)
	}
	return nil, nil
}

func fetchItem[I interface{}](url string) (*I, error) {
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
	var data I
	err = json.Unmarshal([]byte(raw), &data)
	if err != nil {
		log.Error().Str("module", "wallapop").Err(err).Msg("error decoding json")
		return nil, err
	}

	return &data, nil
}
