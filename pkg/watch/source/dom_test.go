package source

import (
	"testing"
)

func TestDOMSource(t *testing.T) {
	source := NewDOMSource("https://www.taylorguitars.jp/events/", ".next_event", nil)
	res, err := source.Fetch()
	if err != nil {
		t.Error(err)
	}
	if res == "" {
		t.Error("cannot be fetched.")
	}
}
