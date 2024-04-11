package pdf

import (
	"io"
	"net/http"
	"os"

	"github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF extracts text from a PDF file at the given URL
func ExtractTextFromPDF(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	f, err := os.CreateTemp("", "paper_*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	content, err := readPDF(f.Name())
	if err != nil {
		return "", err
	}

	return content, nil
}

func readPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		f.Close()
	}()

	var content string
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		content += text
	}
	return content, nil
}
