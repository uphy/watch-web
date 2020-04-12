package check

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding"

	"github.com/PuerkitoBio/goquery"
)

type (
	DOMSource struct {
		URL      string
		Selector string
		Encoding encoding.Encoding
	}
)

func NewDOMSource(url, selector string, encoding encoding.Encoding) *DOMSource {
	return &DOMSource{
		URL:      url,
		Selector: selector,
		Encoding: encoding,
	}
}

func (d *DOMSource) Fetch() (string, error) {
	resp, err := http.Get(d.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}
	var reader io.Reader
	reader = resp.Body
	if d.Encoding != nil {
		reader = d.Encoding.NewDecoder().Reader(reader)
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	doc.Find(d.Selector).Each(func(i int, s *goquery.Selection) {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		text := s.Text()
		text = strings.TrimSpace(text)

		buf.WriteString(text)
	})
	return buf.String(), nil
}
