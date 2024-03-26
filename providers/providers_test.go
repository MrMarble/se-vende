package providers_test

import (
	"testing"

	"gihub.com/mrmarble/se-vende/providers"
)

func TestGetURL(t *testing.T) {
	cases := []struct {
		text     string
		expected string
	}{
		{"", ""},
		{"https://wallapop.com/item/1234", "https://wallapop.com/item/1234"},
		{"https://vinted.es/items/1234", "https://vinted.es/items/1234"},
		{"https://wallapop.com/item/1234?foo=bar", "https://wallapop.com/item/1234?foo=bar"},
		{"https://vinted.es/items/1234?foo=bar", "https://vinted.es/items/1234?foo=bar"},
		{"https://www.wallapop.com/item/1234", "https://www.wallapop.com/item/1234"},
		{"https://es.wallapop.com/item/1234", "https://es.wallapop.com/item/1234"},
		{"https://www.vinted.com/items/1234", "https://www.vinted.com/items/1234"},
		{"this text is ignored https://www.vinted.es/items/245627528-cable-movil?referrer=catalog", "https://www.vinted.es/items/245627528-cable-movil?referrer=catalog"},
		{"google.com", ""},
	}

	for _, c := range cases {
		got := providers.GetURL(c.text)
		if got != c.expected {
			t.Errorf("GetURL(%s) == %s, want %s", c.text, got, c.expected)
		}
	}
}
