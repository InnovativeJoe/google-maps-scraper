package web

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/gosom/google-maps-scraper/deduper"
)

// LoadCSVDeduper reads an existing CSV file and adds all seen links to a new Deduper.
func LoadCSVDeduper(ctx context.Context, csvPath string) (deduper.Deduper, error) {
	d := deduper.New()
	file, err := os.Open(csvPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file does not exist, return an empty Deduper
			return d, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read header row
	headers, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return d, nil
		}
		return nil, err
	}

	linkIndex := -1
	for i, h := range headers {
		if h == "link" {
			linkIndex = i
			break
		}
	}

	// Fallback to index 1 if not found
	if linkIndex == -1 {
		linkIndex = 1
	}

	for {
		row, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if linkIndex < len(row) {
			link := row[linkIndex]
			if link != "" {
				_ = d.AddIfNotExists(ctx, link)
			}
		}
	}

	return d, nil
}
