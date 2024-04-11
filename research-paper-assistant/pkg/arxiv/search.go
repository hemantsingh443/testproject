package arxiv

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Paper represents a research paper with title, authors, summary, and PDF URL
type Paper struct {
	Title    string   `xml:"title"`
	Authors  []string `xml:"author>name"`
	Summary  string   `xml:"summary"`
	PdfURL   string   `xml:"link[title='pdf'] href"`
	ID       string   `xml:"id"`
	Category string   `xml:"category term"`
}

// Feed represents the arXiv API response structure
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Paper  `xml:"entry"`
}

func SearchArxiv(query string, start, maxResults int) ([]Paper, error) {
	url := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=%s&start=%d&max_results=%d", query, start, maxResults)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the XML response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the XML data
	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return feed.Entries, nil
}
